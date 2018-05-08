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

package utilities

import (
	"github.com/jwells131313/goethe"
	"sync"
)

type boundedErrorQueue struct {
	mux sync.Mutex

	capacity uint32
	queue    []goethe.ErrorInformation
}

// NewBoundedErrorQueue creates a new error queue with the given capacity
func NewBoundedErrorQueue(userCapacity uint32) goethe.ErrorQueue {
	return &boundedErrorQueue{
		capacity: userCapacity,
		queue:    make([]goethe.ErrorInformation, 0),
	}
}

func (errorq *boundedErrorQueue) Enqueue(info goethe.ErrorInformation) error {
	if info == nil {
		return nil
	}

	errorq.mux.Lock()
	defer errorq.mux.Unlock()

	if uint32(len(errorq.queue)) >= errorq.capacity {
		return goethe.ErrAtCapacity
	}

	errorq.queue = append(errorq.queue, info)

	return nil
}

func (errorq *boundedErrorQueue) Dequeue() (goethe.ErrorInformation, bool) {
	errorq.mux.Lock()
	defer errorq.mux.Unlock()

	if len(errorq.queue) <= 0 {
		return nil, false
	}

	retVal := errorq.queue[0]
	errorq.queue = errorq.queue[1:]

	return retVal, true
}

func (errorq *boundedErrorQueue) GetSize() int {
	errorq.mux.Lock()
	defer errorq.mux.Unlock()

	return len(errorq.queue)
}

func (errorq *boundedErrorQueue) IsEmpty() bool {
	return errorq.GetSize() == 0
}