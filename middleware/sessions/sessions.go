/*
 * Flow Playground
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sessions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/dapperlabs/flow-playground-api/middleware/httpcontext"
)

type sessionCtxKey string

const sessionCtxKeySession sessionCtxKey = "session"

var store *sessions.CookieStore

// Middleware injects middleware for managing sessions into an HTTP handler.
//
// Sessions are stored using the provided sessions.Store instance.
func Middleware(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Sessions Middleware(): Setting sessionCtxKeySession to", store)
			ctx := context.WithValue(r.Context(), sessionCtxKeySession, store)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func SetCookieStore(cookieStore *sessions.CookieStore) {
	store = cookieStore
}

// Get returns the session with the given name, or creates one if it does not exist.
func Get(ctx context.Context, name string) *sessions.Session {
	// TODO: Store panics on dereference
	// TODO: This means that ctx.Value(sessionCtxKeySession) is not a sessions.Store
	//store := ctx.Value(sessionCtxKeySession).(sessions.Store)
	//fmt.Println("Session.Get(): Store ctx value: ", ctx.Value(sessionCtxKeySession))

	// ignore error because a session is always returned even if one does not exist
	session, _ := store.Get(httpcontext.Request(ctx), name)

	// TODO: Save must be called to set the response cookie/ headers?
	_ = Save(ctx, session) // Pre save in case it's not saved elsewhere?

	return session
}

// Save saves a session by writing it to the HTTP response.
func Save(ctx context.Context, session *sessions.Session) error {
	err := session.Save(
		httpcontext.Request(ctx),
		httpcontext.Writer(ctx),
	)
	if err != nil {
		fmt.Println("Sessions Save(): failed to save session:", err.Error())
		return err
	}

	return nil
}
