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

// Append Handlers to Current List
func (r *ProcessorGroup) Append(handlers ...ProcessHandler) GINGroupProcessor {
	// Add Handlers to List
	r.Chain = append(r.Chain, handlers...)
	return r
}

// Continue Processing the Request (Go to Next Handler)
func (r *ProcessorGroup) Run() GINRequestProcessor {
	// Is Request Processing Finished
	if r.finished {
		panic("ERROR [ GROUP ] EORS Passed")
	}

	// Get GIN Context
	gc := r.GinContext()

	// Do we have Processor Initializer?
	if r.Initializer != nil { // YES: Call Initializer
		handler := r.Initializer
		handler(r, gc)
	}

	// Process Chain
	for !r.finished && r.currentHandler < len(r.Chain) {
		// Call Next Handler
		handler := r.Chain[r.currentHandler]
		handler(r, gc)

		// Skip to Next Handler
		r.currentHandler++
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

func (r *ProcessorGroup) SkipNext() {
	r.currentHandler++
}

func (r *ProcessorGroup) SkipToEnd() {
	// Finish Processing the
	r.currentHandler = len(r.Chain)
	r.finished = true
}
