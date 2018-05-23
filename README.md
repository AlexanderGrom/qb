# qb
[![Go Report Card](https://goreportcard.com/badge/github.com/alexandergrom/qb)](https://goreportcard.com/report/github.com/alexandergrom/qb) [![GoDoc](https://godoc.org/github.com/alexandergrom/qb?status.svg)](https://godoc.org/github.com/alexandergrom/qb)

Qb is a simple query builder.

## Get the package
```bash
$ go get -u github.com/alexandergrom/qb
```

## Examples

Select ...
```go
// Set database grammar (default postgres)
// Not concurrency. You can set it once in the settings
qb.DefaultGrammar("postgres")

// ...

b := new(qb.WhereGroup).
    Where("type", "=", "a").
    WhereOr("type", "=", "b")

q := qb.Query("SELECT id FROM table WHERE %s LIMIT %p OFFSET %p", b, 10, 0)

// SELECT id FROM table WHERE "type" = $1 OR "type" = $2 LIMIT $3 OFFSET $4
fmt.Println(q)

// ["a", "b", 10, 0]
fmt.Println(q.Params())
```

Update ...
```go
b := new(qb.SetGroup).
    Set("name", "Marty").
    Set("surname", "McFly")

q := qb.Query("UPDATE table SET %s WHERE id = %p", b, 10)

// UPDATE table SET "name" = $1, "surname" = $2 WHERE id = $3
fmt.Println(q)

// ["Marty", "McFly", 10]
fmt.Println(q.Params())
```

Insert ...
```go
b := new(qb.ValuesGroup).
    Values(1, "Marty", "McFly")

q := qb.Query("INSERT INTO table (id, name, surname) VALUES %s", b)

// INSERT INTO table (id, name, surname) VALUES ($1, $2, $3)
fmt.Println(q)

// [1, "Marty", "McFly"]
fmt.Println(q.Params())
```

### A more complex example

```go
type (
    CarFilter struct {
        Mark      string
        Model     string
        Color     []int
        Price     []int
        CreatedAt []string
        Limit     int
        Offset    int
    }

    Car struct {
        Mark      string
        Model     string
        Color     int
        Price     int
        CreatedAt time.Time
        UpdatedAt time.Time
    }

    CarRepository struct {
        db *sql.DB
    }
)

func (r *CarRepository) GetByFilter(filter CarFilter) (_ []Car, err error) {
    var builder = new(qb.WhereGroup).WhereRaw("1=1")

    if len(filter.Mark) > 0 {
        builder.Where("mark", "=", filter.Mark)
    }

    if len(filter.Model) > 0 {
        builder.Where("model", "=", filter.Model)
    }

    if len(filter.Color) > 0 {
        var colors = make([]interface{}, len(filter.Color))
        for i, v := range filter.Color {
            colors[i] = v
        }
        builder.WhereIn("color", colors...)
    }

    if len(filter.Price) == 1 {
        builder.Where("price", "<=", filter.Price[0])
    } else if len(filter.Price) == 2 {
        builder.
            Where("price", ">=", filter.Price[0]).
            Where("price", "<=", filter.Price[1])
    }

    if len(filter.CreatedAt) == 1 {
        builder.Where("created_at::data", "=", filter.CreatedAt[0])
    } else if len(filter.CreatedAt) == 2 {
        builder.
            Where("created_at::data", ">=", filter.CreatedAt[0]).
            Where("created_at::data", "<", filter.CreatedAt[1])
    }

    var query = qb.Query(`
        SELECT
            mark, model, color, price, created_at, updated_at
        FROM cars
        WHERE %s
        LIMIT %p
        OFFSET %p
    `, builder, filter.Limit, filter.Offset)

    var rows *sql.Rows
    if rows, err = r.db.Query(query.String(), query.Params()...); err != nil {
        return nil, err
    }
    defer rows.Close()

    cars = make([]Car, 0, filter.Limit)
    for rows.Next() {
        var car = Car{}

        if err = rows.Scan(&car.Mark, &car.Model, &car.Color, &car.Price, &car.CreatedAt, &car.UpdatedAt); err != nil {
            return nil, err
        }

        cars = append(cars, car)
    }

    return cars, rows.Err()
}
```
