/*
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS HEADER.
 *
 * Copyright (c) 2018 Oracle and/or its affiliates. All rights reserved.
 *
 * The contents of this file are subject to the terms of either the GNU
 * General Public License Version 2 only ("GPL") or the Common Development
 * and Distribution License("CDDL") (collectively, the "License").  You
 * may not use this file except in compliance with the License.  You can
 * obtain a copy of the License at
 * https://glassfish.dev.java.net/public/CDDL+GPL_1_1.html
 * or packager/legal/LICENSE.txt.  See the License for the specific
 * language governing permissions and limitations under the License.
 *
 * When distributing the software, include this License Header Notice in each
 * file and include the License file at packager/legal/LICENSE.txt.
 *
 * GPL Classpath Exception:
 * Oracle designates this particular file as subject to the "Classpath"
 * exception as provided by Oracle in the GPL Version 2 section of the License
 * file that accompanied this code.
 *
 * Modifications:
 * If applicable, add the following below the License Header, with the fields
 * enclosed by brackets [] replaced by your own identifying information:
 * "Portions Copyright [year] [name of copyright owner]"
 *
 * Contributor(s):
 * If you wish your version of this file to be governed by only the CDDL or
 * only the GPL Version 2, indicate your decision by adding "[Contributor]
 * elects to include this software in this distribution under the [CDDL or GPL
 * Version 2] license."  If you don't indicate a single choice of license, a
 * recipient has the option to distribute your version of this file under
 * either the CDDL, the GPL Version 2 or to extend the choice of license to
 * its licensees as provided above.  However, if you add GPL Version 2 code
 * and therefore, elected the GPL Version 2 license, then the option applies
 * only if the new code is made subject to such option by the copyright
 * holder.
 */

package goethe

import (
	"sync"
)

// BoundedErrorQueue provides an error queue with a maximum size
type BoundedErrorQueue struct {
	mux sync.Mutex

	capacity uint32
	queue    []ErrorInformation
}

// NewBoundedErrorQueue creates a new error queue with the given capacity
func NewBoundedErrorQueue(userCapacity uint32) ErrorQueue {
	return &BoundedErrorQueue{
		capacity: userCapacity,
		queue:    make([]ErrorInformation, 0),
	}
}

// Enqueue adds an error to the error queue.  If the queue is
// at capacity should return ErrAtCapacity.  All other errors
// will be ignored
func (errorq *BoundedErrorQueue) Enqueue(info ErrorInformation) error {
	if info == nil {
		return nil
	}

	errorq.mux.Lock()
	defer errorq.mux.Unlock()

	if uint32(len(errorq.queue)) >= errorq.capacity {
		return ErrAtCapacity
	}

	errorq.queue = append(errorq.queue, info)

	return nil
}

// Dequeue removes ErrorInformation from the pools
// error queue.  If there were no errors on the queue
// the second return value is false
func (errorq *BoundedErrorQueue) Dequeue() (ErrorInformation, bool) {
	errorq.mux.Lock()
	defer errorq.mux.Unlock()

	if len(errorq.queue) <= 0 {
		return nil, false
	}

	retVal := errorq.queue[0]
	errorq.queue = errorq.queue[1:]

	return retVal, true
}

// GetSize returns the number of items currently in the queue
func (errorq *BoundedErrorQueue) GetSize() int {
	errorq.mux.Lock()
	defer errorq.mux.Unlock()

	return len(errorq.queue)
}

// IsEmpty Returns true if this queue is currently empty
func (errorq *BoundedErrorQueue) IsEmpty() bool {
	return errorq.GetSize() == 0
}
