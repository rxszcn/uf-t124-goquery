package goquery

import "golang.org/x/net/html"

type siblingType int

// Sibling type, used internally when iterating over children at the same
// level (siblings) to specify which nodes are requested.
const (
	siblingPrevUntil siblingType = iota - 3
	siblingPrevAll
	siblingPrev
	siblingAll
	siblingNext
	siblingNextAll
	siblingNextUntil
	siblingAllIncludingNonElements
)

// Find gets the descendants of each element in the current set of matched
// elements, filtered by a selector. It returns a new Selection object
// containing these matched elements.
//
// Note that as for all methods accepting a selector string, the selector is
// compiled and applied by the cascadia package and inherits its behavior and
// constraints regarding supported selectors. See the note on cascadia in
// the goquery documentation here:
// https://github.com/PuerkitoBio/goquery?tab=readme-ov-file#api
func (s *Selection) Find(selector string) *Selection {
	return pushStack(s, findWithMatcher(s.Nodes, compileMatcher(selector)))
}

// FindMatcher gets the descendants of each element in the current set of matched
// elements, filtered by the matcher. It returns a new Selection object
// containing these matched elements.
func (s *Selection) FindMatcher(m Matcher) *Selection {
	return pushStack(s, findWithMatcher(s.Nodes, m))
}

// FindSelection gets the descendants of each element in the current
// Selection, filtered by a Selection. It returns a new Selection object
// containing these matched elements.
func (s *Selection) FindSelection(sel *Selection) *Selection {
	if sel == nil {
		return pushStack(s, nil)
	}
	return s.FindNodes(sel.Nodes...)
}

// FindNodes gets the descendants of each element in the current
// Selection, filtered by some nodes. It returns a new Selection object
// containing these matched elements.
func (s *Selection) FindNodes(nodes ...*html.Node) *Selection {
	return pushStack(s, mapNodes(nodes, func(i int, n *html.Node) []*html.Node {
		if sliceContains(s.Nodes, n) {
			return []*html.Node{n}
		}
		return nil
	}))
}

// Contents gets the children of each element in the Selection,
// including text and comment nodes. It returns a new Selection object
// containing these elements.
func (s *Selection) Contents() *Selection {
	return pushStack(s, getChildrenNodes(s.Nodes, siblingAllIncludingNonElements))
}

// ContentsFiltered gets the children of each element in the Selection,
// filtered by the specified selector. It returns a new Selection
// object containing these elements. Since selectors only act on Element nodes,
// this function is an alias to ChildrenFiltered unless the selector is empty,
// in which case it is an alias to Contents.
func (s *Selection) ContentsFiltered(selector string) *Selection {
	if selector != "" {
		return s.ChildrenFiltered(selector)
	}
	return s.Contents()
}

// ContentsMatcher gets the children of each element in the Selection,
// filtered by the specified matcher. It returns a new Selection
// object containing these elements. Since matchers only act on Element nodes,
// this function is an alias to ChildrenMatcher.
func (s *Selection) ContentsMatcher(m Matcher) *Selection {
	return s.ChildrenMatcher(m)
}

// Children gets the child elements of each element in the Selection.
// It returns a new Selection object containing these elements.
func (s *Selection) Children() *Selection {
	return pushStack(s, getChildrenNodes(s.Nodes, siblingAll))
}

// ChildrenFiltered gets the child elements of each element in the Selection,
// filtered by the specified selector. It returns a new
// Selection object containing these elements.
func (s *Selection) ChildrenFiltered(selector string) *Selection {
	return filterAndPush(s, getChildrenNodes(s.Nodes, siblingAll), compileMatcher(selector))
}

// ChildrenMatcher gets the child elements of each element in the Selection,
// filtered by the specified matcher. It returns a new
// Selection object containing these elements.
func (s *Selection) ChildrenMatcher(m Matcher) *Selection {
	return filterAndPush(s, getChildrenNodes(s.Nodes, siblingAll), m)
}

// Parent gets the parent of each element in the Selection. It returns a
// new Selection object containing the matched elements.
func (s *Selection) Parent() *Selection {
	return pushStack(s, getParentNodes(s.Nodes))
}

// ParentFiltered gets the parent of each element in the Selection filtered by a
// selector. It returns a new Selection object containing the matched elements.
func (s *Selection) ParentFiltered(selector string) *Selection {
	return filterAndPush(s, getParentNodes(s.Nodes), compileMatcher(selector))
}

// ParentMatcher gets the parent of each element in the Selection filtered by a
// matcher. It returns a new Selection object containing the matched elements.
func (s *Selection) ParentMatcher(m Matcher) *Selection {
	return filterAndPush(s, getParentNodes(s.Nodes), m)
}

// Closest gets the first element that matches the selector by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (s *Selection) Closest(selector string) *Selection {
	cs := compileMatcher(selector)
	return s.ClosestMatcher(cs)
}

// ClosestMatcher gets the first element that matches the matcher by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (s *Selection) ClosestMatcher(m Matcher) *Selection {
	return pushStack(s, mapNodes(s.Nodes, func(i int, n *html.Node) []*html.Node {
		// For each node in the selection, test the node itself, then each parent
		// until a match is found.
		for ; n != nil; n = n.Parent {
			if m.Match(n) {
				return []*html.Node{n}
			}
		}
		return nil
	}))
}

// ClosestNodes gets the first element that matches one of the nodes by testing the
// element itself and traversing up through its ancestors in the DOM tree.
func (s *Selection) ClosestNodes(nodes ...*html.Node) *Selection {
	set := make(map[*html.Node]bool, len(nodes))
	for _, n := range nodes {
		set[n] = true
	}
	return pushStack(s, mapNodes(s.Nodes, func(i int, n *html.Node) []*html.Node {
		// For each node in the selection, test the node itself, then each parent
		// until a match is found.
		for ; n != nil; n = n.Parent {
			if set[n] {
				return []*html.Node{n}
			}
		}
		return nil
	}))
}

// ClosestSelection gets the first element that matches one of the nodes in the
// Selection by testing the element itself and traversing up through its ancestors
// in the DOM tree.
func (s *Selection) ClosestSelection(sel *Selection) *Selection {
	if sel == nil {
		return pushStack(s, nil)
	}
	return s.ClosestNodes(sel.Nodes...)
}

// Parents gets the ancestors of each element in the current Selection. It
// returns a new Selection object with the matched elements.
func (s *Selection) Parents() *Selection {
	return pushStack(s, getParentsNodes(s.Nodes, nil, nil))
}

// ParentsFiltered gets the ancestors of each element in the current
// Selection. It returns a new Selection object with the matched elements.
func (s *Selection) ParentsFiltered(selector string) *Selection {
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nil), compileMatcher(selector))
}

// ParentsMatcher gets the ancestors of each element in the current
// Selection. It returns a new Selection object with the matched elements.
func (s *Selection) ParentsMatcher(m Matcher) *Selection {
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nil), m)
}

// ParentsUntil gets the ancestors of each element in the Selection, up to but
// not including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (s *Selection) ParentsUntil(selector string) *Selection {
	return pushStack(s, getParentsNodes(s.Nodes, compileMatcher(selector), nil))
}

// ParentsUntilMatcher gets the ancestors of each element in the Selection, up to but
// not including the element matched by the matcher. It returns a new Selection
// object containing the matched elements.
func (s *Selection) ParentsUntilMatcher(m Matcher) *Selection {
	return pushStack(s, getParentsNodes(s.Nodes, m, nil))
}

// ParentsUntilSelection gets the ancestors of each element in the Selection,
// up to but not including the elements in the specified Selection. It returns a
// new Selection object containing the matched elements.
func (s *Selection) ParentsUntilSelection(sel *Selection) *Selection {
	if sel == nil {
		return s.Parents()
	}
	return s.ParentsUntilNodes(sel.Nodes...)
}

// ParentsUntilNodes gets the ancestors of each element in the Selection,
// up to but not including the specified nodes. It returns a
// new Selection object containing the matched elements.
func (s *Selection) ParentsUntilNodes(nodes ...*html.Node) *Selection {
	return pushStack(s, getParentsNodes(s.Nodes, nil, nodes))
}

// ParentsFilteredUntil is like ParentsUntil, with the option to filter the
// results based on a selector string. It returns a new Selection
// object containing the matched elements.
func (s *Selection) ParentsFilteredUntil(filterSelector, untilSelector string) *Selection {
	return filterAndPush(s, getParentsNodes(s.Nodes, compileMatcher(untilSelector), nil), compileMatcher(filterSelector))
}

// ParentsFilteredUntilMatcher is like ParentsUntilMatcher, with the option to filter the
// results based on a matcher. It returns a new Selection object containing the matched elements.
func (s *Selection) ParentsFilteredUntilMatcher(filter, until Matcher) *Selection {
	return filterAndPush(s, getParentsNodes(s.Nodes, until, nil), filter)
}

// ParentsFilteredUntilSelection is like ParentsUntilSelection, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsFilteredUntilSelection(filterSelector string, sel *Selection) *Selection {
	return s.ParentsMatcherUntilSelection(compileMatcher(filterSelector), sel)
}

// ParentsMatcherUntilSelection is like ParentsUntilSelection, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsMatcherUntilSelection(filter Matcher, sel *Selection) *Selection {
	if sel == nil {
		return s.ParentsMatcher(filter)
	}
	return s.ParentsMatcherUntilNodes(filter, sel.Nodes...)
}

// ParentsFilteredUntilNodes is like ParentsUntilNodes, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection {
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nodes), compileMatcher(filterSelector))
}

// ParentsMatcherUntilNodes is like ParentsUntilNodes, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) ParentsMatcherUntilNodes(filter Matcher, nodes ...*html.Node) *Selection {
	return filterAndPush(s, getParentsNodes(s.Nodes, nil, nodes), filter)
}

// Siblings gets the siblings of each element in the Selection. It returns
// a new Selection object containing the matched elements.
func (s *Selection) Siblings() *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingAll, nil, nil))
}

// SiblingsFiltered gets the siblings of each element in the Selection
// filtered by a selector. It returns a new Selection object containing the
// matched elements.
func (s *Selection) SiblingsFiltered(selector string) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingAll, nil, nil), compileMatcher(selector))
}

// SiblingsMatcher gets the siblings of each element in the Selection
// filtered by a matcher. It returns a new Selection object containing the
// matched elements.
func (s *Selection) SiblingsMatcher(m Matcher) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingAll, nil, nil), m)
}

// Next gets the immediately following sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) Next() *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNext, nil, nil))
}

// NextFiltered gets the immediately following sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextFiltered(selector string) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNext, nil, nil), compileMatcher(selector))
}

// NextMatcher gets the immediately following sibling of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextMatcher(m Matcher) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNext, nil, nil), m)
}

// NextAll gets all the following siblings of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) NextAll() *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextAll, nil, nil))
}

// NextAllFiltered gets all the following siblings of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextAllFiltered(selector string) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextAll, nil, nil), compileMatcher(selector))
}

// NextAllMatcher gets all the following siblings of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) NextAllMatcher(m Matcher) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextAll, nil, nil), m)
}

// Prev gets the immediately preceding sibling of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) Prev() *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrev, nil, nil))
}

// PrevFiltered gets the immediately preceding sibling of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevFiltered(selector string) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrev, nil, nil), compileMatcher(selector))
}

// PrevMatcher gets the immediately preceding sibling of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevMatcher(m Matcher) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrev, nil, nil), m)
}

// PrevAll gets all the preceding siblings of each element in the
// Selection. It returns a new Selection object containing the matched elements.
func (s *Selection) PrevAll() *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevAll, nil, nil))
}

// PrevAllFiltered gets all the preceding siblings of each element in the
// Selection filtered by a selector. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevAllFiltered(selector string) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevAll, nil, nil), compileMatcher(selector))
}

// PrevAllMatcher gets all the preceding siblings of each element in the
// Selection filtered by a matcher. It returns a new Selection object
// containing the matched elements.
func (s *Selection) PrevAllMatcher(m Matcher) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevAll, nil, nil), m)
}

// NextUntil gets all following siblings of each element up to but not
// including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntil(selector string) *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		compileMatcher(selector), nil))
}

// NextUntilMatcher gets all following siblings of each element up to but not
// including the element matched by the matcher. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntilMatcher(m Matcher) *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		m, nil))
}

// NextUntilSelection gets all following siblings of each element up to but not
// including the element matched by the Selection. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntilSelection(sel *Selection) *Selection {
	if sel == nil {
		return s.NextAll()
	}
	return s.NextUntilNodes(sel.Nodes...)
}

// NextUntilNodes gets all following siblings of each element up to but not
// including the element matched by the nodes. It returns a new Selection
// object containing the matched elements.
func (s *Selection) NextUntilNodes(nodes ...*html.Node) *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		nil, nodes))
}

// PrevUntil gets all preceding siblings of each element up to but not
// including the element matched by the selector. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntil(selector string) *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		compileMatcher(selector), nil))
}

// PrevUntilMatcher gets all preceding siblings of each element up to but not
// including the element matched by the matcher. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntilMatcher(m Matcher) *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		m, nil))
}

// PrevUntilSelection gets all preceding siblings of each element up to but not
// including the element matched by the Selection. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntilSelection(sel *Selection) *Selection {
	if sel == nil {
		return s.PrevAll()
	}
	return s.PrevUntilNodes(sel.Nodes...)
}

// PrevUntilNodes gets all preceding siblings of each element up to but not
// including the element matched by the nodes. It returns a new Selection
// object containing the matched elements.
func (s *Selection) PrevUntilNodes(nodes ...*html.Node) *Selection {
	return pushStack(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		nil, nodes))
}

// NextFilteredUntil is like NextUntil, with the option to filter
// the results based on a selector string.
// It returns a new Selection object containing the matched elements.
func (s *Selection) NextFilteredUntil(filterSelector, untilSelector string) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		compileMatcher(untilSelector), nil), compileMatcher(filterSelector))
}

// NextFilteredUntilMatcher is like NextUntilMatcher, with the option to filter
// the results based on a matcher.
// It returns a new Selection object containing the matched elements.
func (s *Selection) NextFilteredUntilMatcher(filter, until Matcher) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		until, nil), filter)
}

// NextFilteredUntilSelection is like NextUntilSelection, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextFilteredUntilSelection(filterSelector string, sel *Selection) *Selection {
	return s.NextMatcherUntilSelection(compileMatcher(filterSelector), sel)
}

// NextMatcherUntilSelection is like NextUntilSelection, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextMatcherUntilSelection(filter Matcher, sel *Selection) *Selection {
	if sel == nil {
		return s.NextMatcher(filter)
	}
	return s.NextMatcherUntilNodes(filter, sel.Nodes...)
}

// NextFilteredUntilNodes is like NextUntilNodes, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		nil, nodes), compileMatcher(filterSelector))
}

// NextMatcherUntilNodes is like NextUntilNodes, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) NextMatcherUntilNodes(filter Matcher, nodes ...*html.Node) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingNextUntil,
		nil, nodes), filter)
}

// PrevFilteredUntil is like PrevUntil, with the option to filter
// the results based on a selector string.
// It returns a new Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntil(filterSelector, untilSelector string) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		compileMatcher(untilSelector), nil), compileMatcher(filterSelector))
}

// PrevFilteredUntilMatcher is like PrevUntilMatcher, with the option to filter
// the results based on a matcher.
// It returns a new Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntilMatcher(filter, until Matcher) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		until, nil), filter)
}

// PrevFilteredUntilSelection is like PrevUntilSelection, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntilSelection(filterSelector string, sel *Selection) *Selection {
	return s.PrevMatcherUntilSelection(compileMatcher(filterSelector), sel)
}

// PrevMatcherUntilSelection is like PrevUntilSelection, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevMatcherUntilSelection(filter Matcher, sel *Selection) *Selection {
	if sel == nil {
		return s.PrevMatcher(filter)
	}
	return s.PrevMatcherUntilNodes(filter, sel.Nodes...)
}

// PrevFilteredUntilNodes is like PrevUntilNodes, with the
// option to filter the results based on a selector string. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevFilteredUntilNodes(filterSelector string, nodes ...*html.Node) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		nil, nodes), compileMatcher(filterSelector))
}

// PrevMatcherUntilNodes is like PrevUntilNodes, with the
// option to filter the results based on a matcher. It returns a new
// Selection object containing the matched elements.
func (s *Selection) PrevMatcherUntilNodes(filter Matcher, nodes ...*html.Node) *Selection {
	return filterAndPush(s, getSiblingNodes(s.Nodes, siblingPrevUntil,
		nil, nodes), filter)
}

// Filter and push filters the nodes based on a matcher, and pushes the results
// on the stack, with the srcSel as previous selection.
func filterAndPush(srcSel *Selection, nodes []*html.Node, m Matcher) *Selection {
	return pushStack(srcSel, winnow(nodes, m, true))
}

// Internal implementation of Find that return raw nodes.
func findWithMatcher(nodes []*html.Node, m Matcher) []*html.Node {
	// Map nodes to find the matches within the children of each node
	return mapNodes(nodes, func(i int, n *html.Node) (result []*html.Node) {
		// Go down one level, becausejQuery's Find selects only within descendants
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				result = append(result, m.MatchAll(c)...)
			}
		}
		return
	})
}

// Internal implementation to get all parent nodes, stopping at the specified
// node (or nil if no stop).
func getParentsNodes(nodes []*html.Node, stopm Matcher, stopNodes []*html.Node) []*html.Node {
	// Collect ancestors inline rather than going through mapNodes, which would
	// allocate a throwaway slice per source node. Ancestor chains of distinct
	// source nodes converge quickly, so deduplicate as we go, exactly like
	// getParentNodes does for the single-parent case. The traversal order is
	// unchanged: source nodes are visited in order, each contributing its
	// ancestors nearest-first, and a node is appended only the first time it
	// is seen.
	var result []*html.Node
	set := make(map[*html.Node]bool, len(nodes))

	for _, n := range nodes {
		for p := n.Parent; p != nil; p = p.Parent {
			if stopm != nil {
				if stopm.Match(p) {
					break
				}
			} else if len(stopNodes) > 0 {
				if isInSlice(stopNodes, p) {
					break
				}
			}
			if p.Type == html.ElementNode && !set[p] {
				set[p] = true
				result = append(result, p)
			}
		}
	}
	return result
}

// Internal implementation of sibling nodes that return a raw slice of matches.
func getSiblingNodes(nodes []*html.Node, st siblingType, untilm Matcher, untilNodes []*html.Node) []*html.Node {
	var f func(*html.Node) bool

	// If the requested siblings are ...Until, create the test function to
	// determine if the until condition is reached (returns true if it is)
	if st == siblingNextUntil || st == siblingPrevUntil {
		f = func(n *html.Node) bool {
			if untilm != nil {
				return untilm.Match(n)
			} else if len(untilNodes) > 0 {
				return isInSlice(untilNodes, n)
			}
			return false
		}
	}

	// Mirror mapNodes' fast paths: with at most one source node there is
	// nothing to deduplicate, so collect directly (the single-node case also
	// keeps the exact-size pre-count inside getChildrenWithSiblingType).
	switch len(nodes) {
	case 0:
		return nil
	case 1:
		return getChildrenWithSiblingType(nodes[0].Parent, st, nodes[0], f)
	}

	// Collect into a single shared slice, deduplicating as we go, instead of
	// building one throwaway slice per source node via mapNodes. Sibling sets
	// of distinct source nodes overlap heavily (siblings of siblings share
	// most of their sibling list), so a set is always used for deduplication,
	// as in getParentNodes.
	var result []*html.Node
	set := make(map[*html.Node]bool, len(nodes))
	for _, n := range nodes {
		result = appendChildrenWithSiblingType(result, set, n.Parent, st, n, f)
	}
	return result
}

// Gets the children nodes of each node in the specified slice of nodes,
// based on the sibling type request.
func getChildrenNodes(nodes []*html.Node, st siblingType) []*html.Node {
	// Same fast paths as getSiblingNodes: with at most one source node there
	// is nothing to deduplicate.
	switch len(nodes) {
	case 0:
		return nil
	case 1:
		return getChildrenWithSiblingType(nodes[0], st, nil, nil)
	}

	// Collect into a single shared slice, deduplicating as we go, instead of
	// building one throwaway slice per source node via mapNodes.
	var result []*html.Node
	set := make(map[*html.Node]bool, len(nodes))
	for _, n := range nodes {
		result = appendChildrenWithSiblingType(result, set, n, st, nil, nil)
	}
	return result
}

// nextSiblingNode advances the sibling walk one step: from cur (nil for the
// first step) it returns the next node to consider for the requested sibling
// type, skipping non-element nodes (except for siblingAllIncludingNonElements)
// and the skipNode. It returns nil when the walk is exhausted. It is a plain
// function rather than a closure so that the per-parent walks below do not
// allocate.
func nextSiblingNode(cur, parent *html.Node, st siblingType, skipNode *html.Node) *html.Node {
	// Based on the sibling type requested, iterate the right way
	for {
		var ret *html.Node
		switch st {
		case siblingAll, siblingAllIncludingNonElements:
			if cur == nil {
				// First iteration, start with first child of parent
				// Skip node if required
				if ret = parent.FirstChild; ret == skipNode && skipNode != nil {
					ret = skipNode.NextSibling
				}
			} else {
				// Skip node if required
				if ret = cur.NextSibling; ret == skipNode && skipNode != nil {
					ret = skipNode.NextSibling
				}
			}
		case siblingPrev, siblingPrevAll, siblingPrevUntil:
			if cur == nil {
				// Start with previous sibling of the skip node
				ret = skipNode.PrevSibling
			} else {
				ret = cur.PrevSibling
			}
		case siblingNext, siblingNextAll, siblingNextUntil:
			if cur == nil {
				// Start with next sibling of the skip node
				ret = skipNode.NextSibling
			} else {
				ret = cur.NextSibling
			}
		default:
			panic("Invalid sibling type.")
		}
		if ret == nil || ret.Type == html.ElementNode || st == siblingAllIncludingNonElements {
			return ret
		}
		// Not a valid node, try again from this one
		cur = ret
	}
}

// appendChildrenWithSiblingType appends the children of parent selected by
// the sibling type st to result and returns the extended slice. When set is
// non-nil, nodes already in the set are skipped and appended nodes are
// recorded in it, so callers collecting the children of several parents into
// one shared slice deduplicate on the fly, without intermediate per-parent
// slices. The traversal order and the until/single-result semantics are
// identical to walking each parent separately and merging the per-parent
// slices with appendWithoutDuplicates.
func appendChildrenWithSiblingType(result []*html.Node, set map[*html.Node]bool, parent *html.Node, st siblingType,
	skipNode *html.Node, untilFunc func(*html.Node) bool) []*html.Node {

	for c := nextSiblingNode(nil, parent, st, skipNode); c != nil; c = nextSiblingNode(c, parent, st, skipNode) {
		// If this is an ...Until case, test before append (returns true
		// if the until condition is reached)
		if st == siblingNextUntil || st == siblingPrevUntil {
			if untilFunc(c) {
				return result
			}
		}
		if set == nil || !set[c] {
			if set != nil {
				set[c] = true
			}
			result = append(result, c)
		}
		if st == siblingNext || st == siblingPrev {
			// Only one node was requested (immediate next or previous), so
			// exit whether or not the candidate was a duplicate
			return result
		}
	}
	return result
}

// Gets the children of the specified parent, based on the requested sibling
// type, skipping a specified node if required.
func getChildrenWithSiblingType(parent *html.Node, st siblingType, skipNode *html.Node,
	untilFunc func(*html.Node) bool) (result []*html.Node) {

	// For the cases that collect every matching sibling, count them in a
	// cheap pointer walk first so the result slice can be sized exactly,
	// avoiding repeated slice growth. The Until cases are skipped (counting
	// would require running the predicate twice) and so are the single-result
	// Next/Prev cases.
	switch st {
	case siblingAll, siblingAllIncludingNonElements, siblingPrevAll, siblingNextAll:
		n := 0
		for c := nextSiblingNode(nil, parent, st, skipNode); c != nil; c = nextSiblingNode(c, parent, st, skipNode) {
			n++
		}
		result = make([]*html.Node, 0, n)
	}

	return appendChildrenWithSiblingType(result, nil, parent, st, skipNode, untilFunc)
}

// Internal implementation of parent nodes that return a raw slice of Nodes.
func getParentNodes(nodes []*html.Node) []*html.Node {
	// Collect parents inline rather than going through mapNodes, which would
	// allocate a throwaway one-element slice per source node. Many source
	// nodes (e.g. siblings) share the same parent, so deduplicate as we go.
	var result []*html.Node
	set := make(map[*html.Node]bool, len(nodes))

	for _, n := range nodes {
		p := n.Parent
		if p == nil || p.Type != html.ElementNode || set[p] {
			continue
		}
		set[p] = true
		result = append(result, p)
	}
	return result
}

// Internal map function used by many traversing methods. Takes the source nodes
// to iterate on and the mapping function that returns an array of nodes.
// Returns an array of nodes mapped by calling the callback function once for
// each node in the source nodes.
func mapNodes(nodes []*html.Node, f func(int, *html.Node) []*html.Node) (result []*html.Node) {
	switch len(nodes) {
	case 0:
		return nil
	case 1:
		return f(0, nodes[0])
	}

	set := make(map[*html.Node]bool, len(nodes))
	for i, n := range nodes {
		if vals := f(i, n); len(vals) > 0 {
			result = appendWithoutDuplicates(result, vals, set)
		}
	}
	return result
}
