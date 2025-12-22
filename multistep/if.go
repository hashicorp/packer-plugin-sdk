// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package multistep

// if returns step only if on is true.
func If(on bool, step Step) Step {
	if !on {
		return &nullStep{}
	}
	return step
}
