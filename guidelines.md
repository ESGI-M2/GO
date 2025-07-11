# Guide de Développement - ORM Go

## Vue d'ensemble
Ce projet implémente un ORM (Object-Relational Mapping) complet en Go, inspiré de Doctrine (PHP), sans dépendances externes. L'ORM est conçu pour être réutilisable dans n'importe quel projet Go.

## Architecture et Principes

### 1. Principes SOLID
- **Single Responsibility Principle (SRP)**: Un fichier = une responsabilité
- **Open/Closed Principle (OCP)**: Ouvert à l'extension, fermé à la modification
- **Liskov Substitution Principle (LSP)**: Les interfaces sont interchangeables
- **Interface Segregation Principle (ISP)**: Interfaces spécifiques et cohérentes
- **Dependency Inversion Principle (DIP)**: Dépendre des abstractions, pas des implémentations

### 2. Structure du Projet
```
project/
├── orm/                    # Cœur de l'ORM
│   ├── core/              # Fonctionnalités principales
│   ├── query/             # Construction de requêtes
│   ├── migration/         # Gestion des migrations
│   └── utils/             # Utilitaires
├── models/                # Modèles de données
├── repository/            # Pattern Repository
├── dialect/              # Support des bases de données
└── tests/                # Tests unitaires et d'intégration
```

### 3. Fonctionnalités Principales

#### 3.1 Gestion des Modèles
- Tags struct pour la configuration (`db`, `primary`, `autoincrement`, `foreign`)
- Mapping automatique Go ↔ SQL
- Support des relations (One-to-One, One-to-Many, Many-to-Many)
- Validation automatique des types

#### 3.2 Construction de Requêtes
- Query Builder fluent API
- Support des clauses WHERE, ORDER BY, LIMIT, OFFSET
- Protection contre les injections SQL
- Optimisation automatique des requêtes

#### 3.3 Gestion des Relations
- Relations automatiques basées sur les tags
- Lazy loading et eager loading
- Cascade operations (CASCADE, SET NULL, RESTRICT)

#### 3.4 Migrations
- Génération automatique des migrations
- Versioning des schémas
- Rollback et rollforward
- Diff detection automatique

#### 3.5 Support Multi-Dialectes
- MySQL (implémenté)
- PostgreSQL (à implémenter)
- SQLite (à implémenter)
- Interface dialecte extensible

### 4. Bonnes Pratiques Go

#### 4.1 Conventions de Nommage
- Variables et fonctions : camelCase
- Types et interfaces : PascalCase
- Constantes : UPPER_SNAKE_CASE
- Packages : lowercase

#### 4.2 Gestion d'Erreurs
- Toujours retourner des erreurs explicites
- Utiliser `errors.Wrap` pour le contexte
- Pas de panics dans le code de production
- Logging structuré

#### 4.3 Performance
- Pool de connexions
- Prepared statements
- Pagination automatique
- Cache des métadonnées

#### 4.4 Tests
- Tests unitaires pour chaque fonction
- Tests d'intégration pour les scénarios complets
- Mocks pour les dépendances externes
- Coverage minimum 80%

### 5. API Design

#### 5.1 Interface Fluent
```go
// Exemple d'utilisation
users, err := orm.Query(&User{}).
    Where("age", ">", 18).
    Where("is_active", "=", true).
    OrderBy("name", "ASC").
    Limit(10).
    Find()
```

#### 5.2 Repository Pattern
```go
type UserRepository struct {
    orm *ORM
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
    // Implémentation
}
```

### 6. Sécurité
- Échappement automatique des paramètres
- Validation des types avant exécution
- Protection contre les injections SQL
- Sanitisation des entrées utilisateur

### 7. Extensibilité
- Interface pour les nouveaux dialectes
- Hooks pour les événements (before/after save, etc.)
- Plugins système
- Configuration flexible

### 8. Documentation
- Documentation Go standard (godoc)
- Exemples d'utilisation
- Guide de migration
- Changelog détaillé

### 9. Monitoring et Debugging
- Logging structuré
- Métriques de performance
- Query logging en mode debug
- Profiling automatique

### 10. Roadmap
- [x] Structure de base
- [x] Query Builder basique
- [ ] Relations avancées
- [ ] Migrations automatiques
- [ ] Support PostgreSQL
- [ ] Cache système
- [ ] Optimisations avancées
- [ ] Documentation complète

## Notes de Développement
- Priorité à la lisibilité et maintenabilité
- Tests obligatoires pour chaque fonctionnalité
- Refactoring continu pour améliorer la qualité
- Documentation à jour avec le code

