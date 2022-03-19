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

func NestedIF(parent GINProcessor, iftest IFProcessHandler, iftrue IFProcessHandler, iffalse IFProcessHandler) *ProcessorIF {
	request := &ProcessorIF{
		IfTest:  iftest,
		IfTrue:  iftrue,
		IfFalse: iffalse,
	}

	// Set Request's Parent
	request.Parent = parent
	return request
}

// Continue Processing the Request (Go to Next Handler)
func (r *ProcessorIF) Run() GINRequestProcessor {
	// Is Request Processing Finished
	if r.finished {
		panic("ERROR [ IF ] EORS Passed")
	}

	// Get GIN Context
	gc := r.GinContext()

	// Do we have Processor Initializer?
	if r.Initializer != nil { // YES: Call Initializer
		handler := r.Initializer
		handler(r, gc)
	}

	// Process Finished?
	if !r.finished { // NO: Continue with if
		// Call Next Handler
		handler := r.IfTest
		handler(*r, gc)
	}

	// Should we call Finalizer?
	if !r.finished && !r.Aborted() && r.Finalizer != nil { // YES: Call
		handler := r.Finalizer
		handler(r, gc)
	}

	// Mark Processor as Finished
	r.finished = true
	return r
}

// Continue Processing with True Handler
func (r *ProcessorIF) ContinueTrue() {
	// Is Request Processing Finished
	if r.finished {
		panic("ERROR [ IF-ContinueTrue ] EORS Passed")
	}

	// Call True Handler
	handler := r.IfTrue
	handler(*r, r.gc)
}

// Continue Processing with True Handler
func (r *ProcessorIF) ContinueFalse() {
	// Is Request Processing Finished
	if r.finished {
		panic("ERROR [ IF-ContinueFalse ] EORS Passed")
	}

	// Call False Handler
	handler := r.IfFalse
	handler(*r, r.gc)
}
