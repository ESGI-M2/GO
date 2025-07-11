# Guidelines de Développement - ORM Go

## Architecture et Principes

### 1. Architecture Modulaire

L'ORM suit une architecture modulaire avec séparation claire des responsabilités :

```
orm/
├── core/           # Fonctionnalités principales
│   ├── interfaces.go  # Définitions des interfaces
│   ├── metadata.go    # Gestion des métadonnées
│   ├── orm.go         # Implémentation principale
│   ├── query_builder.go # Construction de requêtes
│   └── repository.go  # Pattern Repository
├── dialect/        # Support des bases de données
│   ├── interface.go # Interface Dialect
│   └── mysql.go     # Implémentation MySQL
└── models/         # Modèles de données
```

### 2. Principes de Conception

- **Séparation des responsabilités** : Chaque composant a une responsabilité unique
- **Interface-based design** : Utilisation d'interfaces pour la flexibilité
- **Dependency injection** : Injection des dépendances via constructeurs
- **Error handling** : Gestion d'erreurs explicite et informative
- **Type safety** : Utilisation maximale du système de types Go

## Système de Tags ORM

### Nouveau Système de Tags

L'ORM utilise un système de tags concis et puissant pour définir les métadonnées des modèles :

#### Tags de Base

```go
type User struct {
    ID       int    `orm:"pk,auto"`           // Clé primaire + auto-incrément
    Name     string `orm:"index"`              // Champ indexé
    Email    string `orm:"unique"`             // Contrainte unique
    Age      int    `orm:"default:18"`         // Valeur par défaut
    IsActive bool   `orm:"default:true"`       // Booléen avec défaut
    Created  string `orm:"column:created_at"`  // Nom de colonne personnalisé
}
```

#### Tags Avancés

```go
type Post struct {
    ID        int    `orm:"pk,auto"`
    Title     string `orm:"column:post_title,index,length:255"`
    Content   string `orm:"column:post_content"`
    UserID    int    `orm:"fk:users.id"`                    // Clé étrangère
    Status    string `orm:"default:draft,nullable"`          // Défaut + nullable
    Tags      string `orm:"column:post_tags,length:500"`
}
```

#### Référence des Tags

| Tag | Description | Exemple |
|-----|-------------|---------|
| `pk` | Clé primaire | `orm:"pk"` |
| `auto` | Auto-incrément | `orm:"auto"` |
| `unique` | Contrainte unique | `orm:"unique"` |
| `index` | Créer un index | `orm:"index"` |
| `nullable` | Autoriser les valeurs NULL | `orm:"nullable"` |
| `column:name` | Nom de colonne personnalisé | `orm:"column:user_name"` |
| `length:n` | Longueur du champ | `orm:"length:255"` |
| `default:value` | Valeur par défaut | `orm:"default:true"` |
| `fk:table.column` | Clé étrangère | `orm:"fk:users.id"` |

### Compatibilité Ascendante

L'ORM maintient une compatibilité complète avec l'ancien système de tags :

```go
// Ancien style (toujours supporté)
type User struct {
    ID       int    `db:"id" primary:"true" autoincrement:"true"`
    Name     string `db:"name" index:"true"`
    Email    string `db:"email" unique:"true"`
}

// Nouveau style (recommandé)
type User struct {
    ID       int    `orm:"pk,auto"`
    Name     string `orm:"index"`
    Email    string `orm:"unique"`
}
```

## Fonctionnalités Principales

### 1. Gestion des Métadonnées

- **Extraction automatique** via reflection
- **Cache des métadonnées** pour les performances
- **Support des relations** (clés étrangères)
- **Validation des types** Go ↔ SQL

### 2. Query Builder Fluent

```go
// Requêtes complexes avec interface fluide
results, err := orm.Query(&User{}).
    Select("name", "email").
    Where("age", ">", 25).
    Where("is_active", "=", true).
    OrderBy("name", "ASC").
    Limit(10).
    Find()

// Requêtes SQL brutes
rawResults, err := orm.Raw("SELECT * FROM users WHERE age > ?", 25).Find()
```

### 3. Pattern Repository

```go
repo := orm.Repository(&User{})

// Opérations CRUD
err := repo.Save(user)
foundUser, err := repo.Find(1)
allUsers, err := repo.FindAll()
err := repo.Delete(user)

// Recherche par critères
users, err := repo.FindBy(map[string]interface{}{
    "is_active": true,
    "age":       30,
})
```

### 4. Support des Transactions

```go
err := orm.Transaction(func(txORM core.ORM) error {
    // Créer un utilisateur
    user := &User{Name: "John", Email: "john@example.com"}
    repo := txORM.Repository(user)
    err := repo.Save(user)
    if err != nil {
        return err
    }
    
    // Créer un post dans la même transaction
    post := &Post{Title: "Hello", UserID: user.ID}
    postRepo := txORM.Repository(post)
    return postRepo.Save(post)
})
```

## Bonnes Pratiques

### 1. Définition des Modèles

```go
// ✅ Bon - Utilisation du nouveau système de tags
type User struct {
    ID       int    `orm:"pk,auto"`
    Name     string `orm:"index"`
    Email    string `orm:"unique"`
    Age      int    `orm:"default:18"`
    IsActive bool   `orm:"default:true"`
}

// ❌ Éviter - Ancien système verbeux
type User struct {
    ID       int    `db:"id" primary:"true" autoincrement:"true"`
    Name     string `db:"name" index:"true"`
    Email    string `db:"email" unique:"true"`
    Age      int    `db:"age" default:"18"`
    IsActive bool   `db:"is_active" default:"true"`
}
```

### 2. Gestion des Erreurs

```go
// ✅ Bon - Gestion explicite des erreurs
err := orm.Connect(config)
if err != nil {
    log.Fatalf("Échec de connexion: %v", err)
}

// ✅ Bon - Vérification des résultats
user, err := repo.Find(1)
if err != nil {
    return fmt.Errorf("erreur lors de la recherche: %w", err)
}
if user == nil {
    return fmt.Errorf("utilisateur non trouvé")
}
```

### 3. Utilisation des Transactions

```go
// ✅ Bon - Gestion des erreurs dans les transactions
err := orm.Transaction(func(txORM core.ORM) error {
    user := &User{Name: "John", Email: "john@example.com"}
    repo := txORM.Repository(user)
    
    if err := repo.Save(user); err != nil {
        return fmt.Errorf("échec de sauvegarde: %w", err)
    }
    
    return nil
})

if err != nil {
    log.Printf("Transaction échouée: %v", err)
}
```

### 4. Performance

```go
// ✅ Bon - Réutilisation des repositories
repo := orm.Repository(&User{})

for i := 0; i < 100; i++ {
    user := &User{Name: fmt.Sprintf("User%d", i)}
    if err := repo.Save(user); err != nil {
        return err
    }
}

// ❌ Éviter - Création répétée de repositories
for i := 0; i < 100; i++ {
    user := &User{Name: fmt.Sprintf("User%d", i)}
    if err := orm.Repository(user).Save(user); err != nil {
        return err
    }
}
```

## API Design

### 1. Interfaces Principales

```go
// Interface Dialect pour les bases de données
type Dialect interface {
    Connect(config ConnectionConfig) error
    Close() error
    Query(query string, args ...interface{}) (*sql.Rows, error)
    Exec(query string, args ...interface{}) (sql.Result, error)
    Begin() (Transaction, error)
}

// Interface ORM principale
type ORM interface {
    Connect(config ConnectionConfig) error
    Close() error
    RegisterModel(model interface{}) error
    Repository(model interface{}) Repository
    Query(model interface{}) QueryBuilder
    Transaction(fn func(ORM) error) error
}
```

### 2. Méthodes Fluent

```go
// Query Builder avec chaînage
query := orm.Query(&User{}).
    Select("name", "email").
    Where("age", ">", 25).
    Where("is_active", "=", true).
    OrderBy("name", "ASC").
    Limit(10)

results, err := query.Find()
```

## Sécurité

### 1. Protection contre les Injections SQL

- **Requêtes préparées** pour toutes les requêtes
- **Échappement automatique** des paramètres
- **Validation des types** avant exécution

### 2. Gestion des Connexions

```go
// Configuration sécurisée
config := core.ConnectionConfig{
    Driver:   "mysql",
    Host:     "localhost",
    Port:     3306,
    Database: "myapp",
    Username: "app_user",
    Password: "secure_password",
    MaxOpenConns: 10,
    MaxIdleConns: 5,
}
```

## Extensibilité

### 1. Ajout de Nouveaux Dialectes

```go
type PostgreSQLDialect struct {
    db *sql.DB
}

func (p *PostgreSQLDialect) Connect(config ConnectionConfig) error {
    // Implémentation spécifique à PostgreSQL
}

func (p *PostgreSQLDialect) Query(query string, args ...interface{}) (*sql.Rows, error) {
    // Implémentation spécifique à PostgreSQL
}
```

### 2. Hooks et Événements

```go
// Interface pour les hooks
type ModelHooks interface {
    BeforeSave() error
    AfterSave() error
    BeforeDelete() error
    AfterDelete() error
}
```

## Documentation

### 1. Commentaires de Code

```go
// User représente un utilisateur dans le système
type User struct {
    ID       int    `orm:"pk,auto"`           // Identifiant unique
    Name     string `orm:"index"`              // Nom de l'utilisateur
    Email    string `orm:"unique"`             // Email unique
    Age      int    `orm:"default:18"`         // Âge avec défaut
    IsActive bool   `orm:"default:true"`       // Statut actif
}
```

### 2. Exemples d'Utilisation

```go
// Exemple complet d'utilisation
func ExampleUsage() {
    // Initialisation
    mysqlDialect := dialect.NewMySQLDialect()
    orm := core.NewORM(mysqlDialect)
    
    config := core.ConnectionConfig{
        Driver:   "mysql",
        Host:     "localhost",
        Port:     3306,
        Database: "myapp",
        Username: "user",
        Password: "password",
    }
    
    if err := orm.Connect(config); err != nil {
        log.Fatal(err)
    }
    defer orm.Close()
    
    // Enregistrement du modèle
    orm.RegisterModel(&User{})
    
    // Utilisation du repository
    repo := orm.Repository(&User{})
    user := &User{Name: "John", Email: "john@example.com"}
    repo.Save(user)
}
```

## Monitoring et Observabilité

### 1. Logging

```go
// Logging des requêtes
type QueryLogger struct {
    logger *log.Logger
}

func (ql *QueryLogger) LogQuery(query string, args []interface{}, duration time.Duration) {
    ql.logger.Printf("Query: %s, Args: %v, Duration: %v", query, args, duration)
}
```

### 2. Métriques

```go
// Métriques de performance
type Metrics struct {
    QueryCount    int64
    QueryDuration time.Duration
    ErrorCount    int64
}
```

## Roadmap

### Phase 1 - Base (✅ Complétée)
- [x] Architecture modulaire
- [x] Système de métadonnées
- [x] Query Builder fluent
- [x] Pattern Repository
- [x] Support MySQL
- [x] Nouveau système de tags ORM

### Phase 2 - Avancé (🔄 En cours)
- [ ] Support PostgreSQL
- [ ] Support SQLite
- [ ] Relations avancées (One-to-Many, Many-to-Many)
- [ ] Système de migrations automatiques
- [ ] Cache et optimisations

### Phase 3 - Production (📋 Planifié)
- [ ] Hooks et événements
- [ ] Validation automatique
- [ ] Code generation
- [ ] Documentation automatique
- [ ] Outils de migration
- [ ] Monitoring avancé

## Tests et Qualité

### 1. Tests Unitaires

```go
func TestUserModel(t *testing.T) {
    mm := NewMetadataManager()
    user := &User{}
    
    metadata, err := mm.ExtractMetadata(user)
    if err != nil {
        t.Fatalf("ExtractMetadata should not return error: %v", err)
    }
    
    if metadata.TableName != "user" {
        t.Errorf("Expected table name 'user', got '%s'", metadata.TableName)
    }
}
```

### 2. Tests d'Intégration

```go
func TestUserCRUD(t *testing.T) {
    // Test complet des opérations CRUD
    orm := setupTestORM(t)
    defer orm.Close()
    
    user := &User{Name: "Test User", Email: "test@example.com"}
    
    // Test Create
    err := orm.Repository(user).Save(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    // Test Read
    found, err := orm.Repository(user).Find(user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Name, found.(*User).Name)
}
```

### 3. Couverture de Code

```bash
# Génération du rapport de couverture
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Conclusion

Ce guide fournit les principes et bonnes pratiques pour développer avec l'ORM Go. Le nouveau système de tags ORM offre une syntaxe plus concise et expressive tout en maintenant la compatibilité avec l'ancien système. L'architecture modulaire permet une extensibilité facile et une maintenance simplifiée.

