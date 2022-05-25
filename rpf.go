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

// TODO: Current Processors don't allow for asynchronous processing!!

type GINRequestProcessor interface {
	IsFinished() bool                                    // Has Processor Terminated?
	Aborted() bool                                       // Was Processor Aborted
	Run() GINRequestProcessor                            // Run Processor
	GinContext() *gin.Context                            // GIN Request Context
	Abort(code int, data *gin.H)                         // Abort GIN Request with Error
	Answer(code int)                                     // Complete GIN Request with Code
	AnswerWithData(code int, data *gin.H)                // Complete GIN Request with Code and Data
	ResponseCode() int                                   // Current Request Response Code
	ResponseData() interface{}                           // Current Request Response Data
	SetResponseCode(code int)                            // SET Request Response Data
	SetResponseData(data *gin.H)                         // SET GIN Response Message
	SetResponseDataValue(name string, value interface{}) // SET Response Message Object Value
}

type GINProcessorContext interface {
	Has(name string) bool
	HasLocal(name string) bool
	Get(name string) interface{}
	MustGet(name string) interface{}
	Set(name string, value interface{}) (interface{}, bool)
	SetLocal(name string, value interface{}) (interface{}, bool)
	Unset(name string) (interface{}, bool)
	LocalToGlobal(name string)
}

type GINGroupProcessor interface {
	SkipNext()
	SkipToEnd()
	Append(handlers ...ProcessHandler) GINGroupProcessor
}

type GINProcessor interface {
	GINRequestProcessor
	GINProcessorContext
}

type BaseProcessor struct {
	GINProcessor
	finished bool                   // Processor Finished?
	gc       *gin.Context           // gin Context
	locals   map[string]interface{} // Local Context Variables
}

// Standard Process Handler Function
type ProcessHandler func(GINProcessor, *gin.Context)

// Chain of Process Functions
type ProcessChain []ProcessHandler

type ChildProcessor struct {
	BaseProcessor
	Parent      GINProcessor   // Parent Processor (if any)
	Initializer ProcessHandler // USE CASE: Function that Creates Common Handler Chain, but, can have different initializer
	Finalizer   ProcessHandler // USE CASE: Function that Creates Common Handler Chain, but, can have different termination
}

// Process Group
type ProcessorGroup struct {
	ChildProcessor
	GINGroupProcessor
	currentHandler int          // Stage Being Processed in Chain
	Chain          ProcessChain // Chain of Handlers
}

// Root Processor Object
type ProcessorRoot struct {
	BaseProcessor
	GINGroupProcessor
	aborted         bool           // FLAG: Process Was Aborted
	title           string         // Text Used to Identify Request Type
	currentHandler  int            // Stage Being Processed in Chain
	Chain           ProcessChain   // Chain of Handlers
	responseCode    int            // API Response Code
	responseData    *gin.H         // API Response Data
	responseHandler ProcessHandler // Handler Used to Generate Answer
}

// IF Process Handler Function
type IFProcessHandler func(ProcessorIF, *gin.Context)

// IF Processor
type ProcessorIF struct {
	ChildProcessor
	IfTest  IFProcessHandler
	IfTrue  IFProcessHandler
	IfFalse IFProcessHandler
}
