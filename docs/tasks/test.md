#app.gtx
```tsx
package main

import "user"

var users = []string{"Alice", "Bob", "Charlie"}

func Render() gtk.Widget {
    return <Box orientation={gtk.ORIENTATION_VERTICAL} spacing={8}>
        {for _, name := range users {
            <user Name={name} />
        }}
    </Box>
}
```

#user.gtx
```tsx
package user

type Props struct {
    Name string
}

func Render(props Props) gtk.Widget {
    return <Box orientation={gtk.ORIENTATION_VERTICAL} spacing={4}>
        <Label label={props.Name} />
    </Box>
}
```