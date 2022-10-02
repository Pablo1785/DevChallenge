package main

import (
	"DevChallenge/handlers"
	"DevChallenge/handlers/middlewares"
	"DevChallenge/infrastructure/dao/impl"
	impl2 "DevChallenge/service/impl"
	"github.com/gorilla/mux"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
	"net/http"
)

func helloWorld(uri, username, password string) (string, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (any, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]any{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}

func main() {
	// Database
	driver, err := neo4j.NewDriver("neo4j://db:7687", neo4j.BasicAuth("neo4j", "s3cr3t", ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return tx.Run("CALL apoc.schema.assert(null, {Person: ['id']}, false)", nil)
	})
	if err != nil {
		panic(err)
	}

	// Data Access Layer
	messageDao := impl.NewMessageDao(&session)
	personDao := impl.NewPersonDao(&session)

	// Service Layer
	messageService := impl2.NewMessageService(messageDao)
	personService := impl2.NewPersonService(personDao)

	// Presentation Layer
	messageServiceHandler := handlers.NewMessageServiceHandler(messageService)
	personServiceHandler := handlers.NewPersonServiceHandler(personService)

	// ROUTES
	r := mux.NewRouter().StrictSlash(true)
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/messages", messageServiceHandler.SendMessage()).Methods(http.MethodPost)
	apiRouter.HandleFunc("/people", personServiceHandler.AddPerson()).Methods(http.MethodPost)
	apiRouter.HandleFunc("/people/{personId}/trust_connections", personServiceHandler.CreateOrUpdateTrustConnections()).Methods(http.MethodPost)

	// MIDDLEWARE
	r.Use(middlewares.CommonResponseHeadersFunc)

	// RUN SERVER
	logrus.Fatal(http.ListenAndServe(":8080", r))
}
