# Go ORM - A Complete Relational ORM

Un ORM (Object-Relational Mapping) complet écrit en Go, inspiré de Doctrine (PHP), sans dépendances externes. Conçu pour être réutilisable dans n'importe quel projet Go.

## 🚀 Fonctionnalités

### ✅ Implémentées
- **Architecture modulaire** avec interfaces claires
- **Gestion des métadonnées** automatique via reflection
- **Query Builder fluent** avec chaînage de méthodes
- **Pattern Repository** pour les opérations CRUD
- **Support des transactions** avec rollback automatique
- **Dialectes de base de données** (MySQL implémenté)
- **Mapping automatique** Go ↔ SQL
- **Support des relations** (clés étrangères)
- **Requêtes SQL brutes** pour les cas complexes
- **Tests unitaires complets**

### 🔄 En cours de développement
- Support PostgreSQL et SQLite
- Relations avancées (One-to-Many, Many-to-Many)
- Système de migrations automatiques
- Cache et optimisations
- Hooks et événements

## 📦 Installation

```bash
go get github.com/votre-username/go-orm
```

## 🏗️ Architecture

```
project/
├── orm/
│   ├── core/              # Interfaces et implémentations principales
│   │   ├── interfaces.go  # Définitions des interfaces
│   │   ├── metadata.go    # Gestion des métadonnées
│   │   ├── orm.go         # Implémentation principale
│   │   ├── query_builder.go # Construction de requêtes
│   │   └── repository.go  # Pattern Repository
│   ├── sql/               # Composants SQL
│   └── utils/             # Utilitaires
├── dialect/               # Support des bases de données
│   └── mysql.go          # Dialecte MySQL
├── models/                # Modèles de données
└── examples/              # Exemples d'utilisation
```

## 🎯 Utilisation Rapide

### 1. Définir vos modèles

```go
type User struct {
    ID        int       `db:"id" primary:"true" autoincrement:"true"`
    Name      string    `db:"name"`
    Email     string    `db:"email" unique:"true"`
    Age       int       `db:"age"`
    IsActive  bool      `db:"is_active"`
    CreatedAt time.Time `db:"created_at"`
}
```

### 2. Initialiser l'ORM

```go
// Créer le dialecte MySQL
mysqlDialect := dialect.NewMySQLDialect()

// Créer l'instance ORM
orm := core.NewORM(mysqlDialect)

// Configurer la connexion
config := core.ConnectionConfig{
    Driver:   "mysql",
    Host:     "localhost",
    Port:     3306,
    Database: "myapp",
    Username: "user",
    Password: "password",
}

// Se connecter
if err := orm.Connect(config); err != nil {
    log.Fatal(err)
}
defer orm.Close()
```

### 3. Enregistrer les modèles

```go
// Enregistrer les modèles
if err := orm.RegisterModel(&User{}); err != nil {
    log.Fatal(err)
}

// Créer les tables
if err := orm.Migrate(); err != nil {
    log.Fatal(err)
}
```

### 4. Utiliser l'ORM

#### Opérations CRUD basiques

```go
// Créer un utilisateur
user := &User{
    Name:      "John Doe",
    Email:     "john@example.com",
    Age:       30,
    IsActive:  true,
    CreatedAt: time.Now(),
}

// Sauvegarder (insert)
if err := orm.Repository(&User{}).Save(user); err != nil {
    log.Fatal(err)
}

// Trouver par ID
foundUser, err := orm.Repository(&User{}).Find(user.ID)
if err != nil {
    log.Fatal(err)
}

// Mettre à jour
user.Age = 31
if err := orm.Repository(&User{}).Update(user); err != nil {
    log.Fatal(err)
}

// Supprimer
if err := orm.Repository(&User{}).Delete(user); err != nil {
    log.Fatal(err)
}
```

#### Query Builder

```go
// Requête avec conditions
users, err := orm.Query(&User{}).
    Where("age", ">", 25).
    Where("is_active", "=", true).
    OrderBy("name", "ASC").
    Limit(10).
    Find()

// Compter les résultats
count, err := orm.Query(&User{}).
    Where("is_active", "=", true).
    Count()

// Vérifier l'existence
exists, err := orm.Query(&User{}).
    Where("email", "=", "john@example.com").
    Exists()
```

#### Repository Pattern

```go
repo := orm.Repository(&User{})

// Trouver par critères
users, err := repo.FindBy(map[string]interface{}{
    "is_active": true,
    "age":       30,
})

// Trouver un seul par critères
user, err := repo.FindOneBy(map[string]interface{}{
    "email": "john@example.com",
})

// Compter tous
count, err := repo.Count()
```

#### Transactions

```go
err := orm.Transaction(func(txORM core.ORM) error {
    // Créer un utilisateur
    user := &User{Name: "Alice", Email: "alice@example.com"}
    if err := txORM.Repository(&User{}).Save(user); err != nil {
        return err
    }
    
    // Créer un post lié à l'utilisateur
    post := &Post{Title: "Hello", UserID: user.ID}
    if err := txORM.Repository(&Post{}).Save(post); err != nil {
        return err
    }
    
    return nil
})
```

#### SQL brut

```go
// Requête SQL brute
results, err := orm.Raw("SELECT COUNT(*) as count FROM users WHERE age > ?", 25).Find()

// Requête complexe
complexResults, err := orm.Raw(`
    SELECT u.name, COUNT(p.id) as post_count 
    FROM users u 
    LEFT JOIN posts p ON u.id = p.user_id 
    WHERE u.is_active = ? 
    GROUP BY u.id, u.name
`, true).Find()
```

## 🏷️ Tags de modèles

### Tags de base

```go
type User struct {
    ID        int    `db:"id" primary:"true" autoincrement:"true"`
    Name      string `db:"name"`
    Email     string `db:"email" unique:"true"`
    Age       int    `db:"age" index:"true"`
    IsActive  bool   `db:"is_active"`
}
```

### Tags disponibles

- `db:"column_name"` - Nom de la colonne en base
- `primary:"true"` - Clé primaire
- `autoincrement:"true"` - Auto-incrémentation
- `unique:"true"` - Contrainte unique
- `index:"true"` - Index sur la colonne
- `length:"255"` - Longueur pour VARCHAR
- `default:"value"` - Valeur par défaut
- `foreign:"table.column"` - Clé étrangère
- `ondelete:"CASCADE"` - Action ON DELETE
- `onupdate:"CASCADE"` - Action ON UPDATE

## 🧪 Tests

```bash
# Lancer tous les tests
go test ./...

# Tests avec couverture
go test -cover ./...

# Tests spécifiques
go test ./orm/core -v
```

## 📚 Exemples

Voir le dossier `examples/` pour des exemples complets d'utilisation.

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit les changements (`git commit -m 'Add some AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🎯 Roadmap

- [x] Architecture de base
- [x] Query Builder
- [x] Repository Pattern
- [x] Transactions
- [x] Support MySQL
- [ ] Support PostgreSQL
- [ ] Support SQLite
- [ ] Relations avancées
- [ ] Système de migrations
- [ ] Cache et optimisations
- [ ] Documentation complète

## 📞 Support

Pour toute question ou problème, veuillez ouvrir une issue sur GitHub.