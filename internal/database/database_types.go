package database

import "time"

type Task struct {
	ID                int        `db:"id" json:"id"`
	Number            int        `db:"number" json:"number"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	SolvedAt          *time.Time `db:"solved_at" json:"solved_at"`
	PlatformDifficult int        `db:"platform_difficult" json:"platform_difficult"`
	MyDifficult       Difficulty `db:"my_difficult" json:"my_difficult"`
	SolvedWithHint    bool       `db:"solved_with_hint" json:"solved_with_hint"`
	Description       string     `db:"description" json:"description"`
	IsMasthaved       bool       `db:"is_masthaved" json:"is_masthaved"`
	Labels            []Label    `db:"labels" json:"labels"`
}

type Difficulty int

const (
	Easy Difficulty = iota + 1
	Medium
	Hard
)

type Label int

const (
	Massive Label = iota
	String
	HashTable
	Math
	DynamicProgramming
	Sorting
	Greedy
	DepthFirstSearch
	BinarySearch
	DataBase
	Matrix
	BitManipulation
	Tree
	BreadthFirstSearch
	TwoPointers
	PrefixSum
	Heap
	Simulation
	Counting
	Graph
	BinaryTree
	Stack
	SlidingWindow
	Design
	Enumeration
	Backtracking
	UnionFind
	NumberTheory
	LinkedList
	OrderedSet
	SegmentTree
	MonotonicStack
	Trie
	DivideAndConquer
	Combinatorics
	Bitmask
	Queue
	Recursion
	Geometry
	BinaryIndexedTree
	Memoization
	HashFunction
	BinarySearchTree
	ShortestPath
	StringMatching
	TopologicalSort
	RollingHash
	GameTheory
	Interactive
	DataStream
	MonotonicQueue
	Brainteaser
	DoubleLinkedList
	MergeSort
	Randomized
	CountingSort
	Iterator
	Concurrency
	SuffixArray
	LineSweep
	ProbabilityAndStatistics
	Quickselect
	MinimumSpanningTree
	BucketSort
	Shell
)
