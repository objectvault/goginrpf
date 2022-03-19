// cSpell:ignore eors, goginrpf, gonic, paulo ferreira, pferreira, testto

// Package GO GIN Request Processing Framework
package goginrpf

/*
 * This file is part of the ObjectVault Project.
 * Copyright (C) 2020-2022 Paulo Ferreira <vault at sourcenotes.org>
 *
 * This work is published under the GNU AGPLv3.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"github.com/gin-gonic/gin"
)

// Has Processor Terminated?
func (r *ChildProcessor) IsFinished() bool {
	return r.finished
}

func (r *ChildProcessor) Aborted() bool {
	return r.Parent.Aborted()
}

// Abort Processing Due To Error
func (r *ChildProcessor) Abort(code int, data *gin.H) {
	r.finished = true
	r.Parent.Abort(code, data)
}

func (r *ChildProcessor) Answer(code int) {
	r.finished = true
	r.Parent.Answer(code)
}

func (r *ChildProcessor) AnswerWithData(code int, data *gin.H) {
	r.finished = true
	r.Parent.AnswerWithData(code, data)
}

func (r *ChildProcessor) GinContext() *gin.Context {
	// Do we have Parent Context Cached?
	if r.gc == nil { // NO: Get Result and Cache
		r.gc = r.Parent.GinContext()
	}
	return r.gc
}

func (r *ChildProcessor) ResponseCode() int {
	return r.Parent.ResponseCode()
}

func (r *ChildProcessor) ResponseData() interface{} {
	return r.Parent.ResponseData()
}

func (r *ChildProcessor) SetResponseData(data *gin.H) {
	// ONLY Process ROOT Maintains Data Value
	r.Parent.SetResponseData(data)
}

func (r *ChildProcessor) SetReponseDataValue(name string, value interface{}) {
	// ONLY Process ROOT Maintains Data Value
	r.Parent.SetReponseDataValue(name, value)
}

func (r *ChildProcessor) Has(name string) bool {
	// Do we have Local Values?
	if !r.HasLocal(name) { // NO: Try Parent
		return r.Parent.Has(name)
	}
	// ELSE: Yes
	return true
}

func (r *ChildProcessor) HasLocal(name string) bool {
	// Do we have Local Values?
	if r.locals != nil { // YES
		// Do we have a Local Value?
		_, exists := r.locals[name]
		return exists
	}
	// ELSE: No Local Value
	return false
}

// Get Context Variable
func (r *ChildProcessor) Get(name string) interface{} {
	// Do we have Local Values?
	if r.locals != nil { // YES
		// Do we have a Local Value?
		value, ok := r.locals[name]
		if ok { // YES
			return value
		}
	}
	// ELSE: Try Parent
	return r.Parent.Get(name)
}

// MustGet Context Variable
func (r *ChildProcessor) MustGet(name string) interface{} {
	// Do we have Local Values?
	if r.locals != nil { // YES
		// Do we have a Local Value?
		value, ok := r.locals[name]
		if ok { // YES
			return value
		}
	}
	// ELSE: Try Parent
	return r.Parent.MustGet(name)
}

// Set Global Context Variable
func (r *ChildProcessor) Set(name string, value interface{}) (interface{}, bool) {
	return r.Parent.Set(name, value)
}

// Set Local Context Variable
func (r *ChildProcessor) SetLocal(name string, value interface{}) (interface{}, bool) {
	var old interface{}
	exists := false
	if r.locals == nil {
		r.locals = make(map[string]interface{})
	} else {
		old, exists = r.locals[name]
	}

	r.locals[name] = value
	return old, exists
}

// Unset Delete Context Variable
func (r *ChildProcessor) Unset(name string) (interface{}, bool) {
	// Do we have local value?
	if r.HasLocal(name) { // YES: Unset that
		old, _ := r.locals[name]
		delete(r.locals, name)
		return old, true
	}
	// ELSE: Unset Parent's Value
	return r.Parent.Unset(name)
}

// Move from Local Context to Global Context
func (r *ChildProcessor) LocalToGlobal(name string) {
	// Do we have local value?
	if r.HasLocal(name) { // YES: Unset that
		v, _ := r.locals[name]
		delete(r.locals, name)
		r.Parent.Set(name, v)
	}
}
