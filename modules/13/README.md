```bash
go run github.com/99designs/gqlgen generate # refaz o código a partir do schema.graphqls
```

```graphql
mutation createCategory {
    createCategory(input: {
        name: "Iniciante"
    }) {
        id
        name
        description
    }
}
```

```graphql
query categories{
    categories{
        id
        name
        description
    }
}
```

```graphql
mutation createCourse {
    createCourse(input: {
        title: "JavaScript Avançado"
        description: "curso iniciante de JS"
        categoryId: "823153af-48ba-40fb-aec5-9842b6a00f22"
    }) {
        id
        title
        description
        category {
            id
        }
    }
}
```

```graphql
query courses {
    courses{
        id
        title
        description
        category {
            id
        }
    }
}
```
