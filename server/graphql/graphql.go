package graphql

import (
	"fmt"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/graphql"
	"github.com/animenotifier/notify.moe/arn"
)

func Install(app *aero.Application) {
	api := graphql.New(arn.DB)

	// Block private collections
	api.AddRootResolver(func(name string, arguments graphql.Map) (interface{}, error, bool) {
		typeName := strings.TrimPrefix(name, "all")
		typeName = strings.TrimPrefix(typeName, "like")

		if arn.IsPrivateType(typeName) {
			return nil, fmt.Errorf("Type '%s' is private", typeName), true
		}

		return nil, nil, false
	})

	// Like objects
	// api.AddRootResolver(func(name string, arguments graphql.Map) (interface{}, error, bool) {
	// 	if !strings.HasPrefix(name, "like") {
	// 		return nil, nil, false
	// 	}

	// 	id, ok := arguments["ID"].(string)

	// 	if !ok {
	// 		return nil, fmt.Errorf("'%s' needs to specify an ID", name), true
	// 	}

	// 	typeName := strings.TrimPrefix(name, "like")
	// 	obj, err := arn.DB.Get(typeName, id)

	// 	if err != nil {
	// 		return nil, err, true
	// 	}

	// 	field := reflect.ValueOf(obj).Elem().FieldByName("IsDraft")

	// 	if field.IsValid() && field.Bool() {
	// 		return nil, errors.New("Drafts need to be published before they can be liked"), true
	// 	}

	// 	likeable, ok := obj.(arn.Likeable)

	// 	if !ok {
	// 		return nil, fmt.Errorf("'%s' does not implement the Likeable interface", name), true
	// 	}

	// 	// TODO: Authentication
	// 	// user := GetUserFromContext(ctx)

	// 	// if user == nil {
	// 	// 	return errors.New("Not logged in")
	// 	// }

	// 	// likeable.Like(user.ID)

	// 	// Call OnLike if the object implements it
	// 	// receiver, ok := likeable.(LikeEventReceiver)

	// 	// if ok {
	// 	// 	receiver.OnLike(user)
	// 	// }

	// 	likeable.Save()
	// 	return obj, nil, true
	// })

	app.Post("/graphql", api.Handler())
}
