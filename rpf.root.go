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

// Constructor Create an RPF Instance
func RootProcessor(title string, c *gin.Context, code int, answer ProcessHandler) *ProcessorRoot {
	request := &ProcessorRoot{
		title:           title,
		responseCode:    code,
		responseHandler: answer,
	}

	// Set GIN Context
	request.gc = c
	return request
}

// Has Processor Terminated?
func (r *ProcessorRoot) IsFinished() bool {
	return r.finished
}

func (r *ProcessorRoot) Aborted() bool {
	return r.aborted
}

// Abort Processing Due To Error
func (r *ProcessorRoot) Abort(code int, data *gin.H) {
	r.AnswerWithData(code, data)
	r.aborted = true
}

func (r *ProcessorRoot) Answer(code int) {
	r.responseCode = code
	r.finished = true
}

func (r *ProcessorRoot) AnswerWithData(code int, data *gin.H) {
	r.responseCode = code
	r.responseData = data
	r.finished = true
}

func (r *ProcessorRoot) GinContext() *gin.Context {
	return r.gc
}

// ResponseCode Current Response Code
func (r *ProcessorRoot) ResponseCode() int {
	return r.responseCode
}

// ResponseData Current Response Data
func (r *ProcessorRoot) ResponseData() interface{} {
	return r.responseData
}

// SetResponseData Set Response Code
func (r *ProcessorRoot) SetResponseCode(code int) {
	// Set New Response Code
	r.responseCode = code
}

// SetResponseData Set Response Data
func (r *ProcessorRoot) SetResponseData(data *gin.H) {
	// Set New Response Data
	r.responseData = data
}

// SetResponseDataValue Set A Property in the Response Data
func (r *ProcessorRoot) SetResponseDataValue(name string, value interface{}) {
	// Do we already have response data?
	if r.responseData == nil { // NO: Initialize
		r.responseData = &gin.H{}
	}

	(*r.responseData)[name] = value
}

// Append Handlers to Current List
func (r *ProcessorRoot) Append(handlers ...ProcessHandler) GINGroupProcessor {
	// Add Handlers to List
	r.Chain = append(r.Chain, handlers...)
	return r
}

// Continue Processing the Request (Go to Next Handler)
func (r *ProcessorRoot) Run() GINRequestProcessor {
	// Is Request Processing Finished
	if r.finished {
		panic("ERROR [ ROOT ] EORS Passed")
	}

	// Process Chain
	for !r.finished && r.currentHandler < len(r.Chain) {
		// Call Next Handler
		handler := r.Chain[r.currentHandler]
		handler(r, r.gc)

		// Skip to Next Handler
		r.currentHandler++
	}

	// Mark Processor as Finished
	r.finished = true
	r.SkipToEnd()
	return r
}

// Skip to the Next Stage
func (r *ProcessorRoot) SkipNext() {
	r.currentHandler++
}

// SkipToAnswer Finish Request
func (r *ProcessorRoot) SkipToEnd() {
	// Finish Processing the
	r.currentHandler = len(r.Chain)
	r.finished = true

	// Call Response Generator
	r.responseHandler(r, r.gc)
}

// Title Request Title
func (r *ProcessorRoot) Title() string {
	return r.title
}

func (r *ProcessorRoot) Has(name string) bool {
	// Alias for HasLocal
	return r.HasLocal(name)
}

func (r *ProcessorRoot) HasLocal(name string) bool {
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
func (r *ProcessorRoot) Get(name string) interface{} {
	// Do we have Local Values?
	if r.locals != nil { // YES
		// Do we have a Local Value?
		value, ok := r.locals[name]
		if ok { // YES
			return value
		}
	}
	// ELSE: Return nil
	return nil
}

// MustGet Context Variable
func (r *ProcessorRoot) MustGet(name string) interface{} {
	// Do we have Local Values?
	if r.locals != nil { // YES
		// Do we have a Local Value?
		value, ok := r.locals[name]
		if ok { // YES
			return value
		}
	}
	// ELSE: Return nil
	panic("Key \"" + name + "\" does not exist")
}

// Set Context Variable
func (r *ProcessorRoot) Set(name string, value interface{}) (interface{}, bool) {
	// Alias for SetLocal
	return r.SetLocal(name, value)
}

// Set Local Context Variable
func (r *ProcessorRoot) SetLocal(name string, value interface{}) (interface{}, bool) {
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
func (r *ProcessorRoot) Unset(name string) (interface{}, bool) {
	// Do we have local value?
	if r.HasLocal(name) { // YES: Unset that
		old, _ := r.locals[name]
		delete(r.locals, name)
		return old, true
	}
	return nil, false
}

// Move from Local Context to Global Context
func (r *ProcessorRoot) LocalToGlobal(name string) {
	// Do Nothing
}
