package iterator

import (
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// Examples

func ExampleFromSlice() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// FromSlice turns the slice into an iterator.
	si := FromSlice(s)

	// Print each value from the slice iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	_ = ForEach[int](si, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
}

func ExampleFromChannel() {
	c := make(chan int)

	go func() {
		defer close(c)
		for i := 1; i <= 10; i++ {
			c <- i
		}
	}()

	// FromChannel turns the channel into an iterator.
	si := FromChannel(c)

	// Print each value from the slice iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	_ = ForEach[int](si, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
}

func ExampleSequence() {
	// Instead of doing the following we can also use generators.
	// Generators are Iterators that make up data while iterating.
	// This example uses the Sequence generator.

	// s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// si := FromSlice(s)

	// Get a sequence iterator that generates values from 1 to 10.
	si := Sequence(1, 10)

	// Print each value from the sequence iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	_ = ForEach[int](si, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// 10
}

func ExampleFilter() {
	// odd is a predicate that evaluates to true when an odd number is encountered.
	odd := func(v int) bool {
		return (v % 2) != 0
	}

	// Get a sequence iterator that generates values from 1 to 10.
	si := Sequence(1, 10)
	// Get a filter iterator that takes the sequence iterator as an input and uses the odd predicate to filter.
	// the values from the sequence iterator.
	fi := Filter[int](si, odd)

	// Print each value from the filter iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	_ = ForEach[int](fi, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 1
	// 3
	// 5
	// 7
	// 9
}

func ExampleMap() {
	// double is a map closure that doubles each value.
	double := func(v int) int {
		return v * 2
	}

	// Get a sequence iterator that generates values from 1 to 10.
	si := Sequence(1, 10)
	// Get a map iterator that takes the sequence iterator as an input and uses the map closure.
	// to double the values from the sequence iterator.
	mi := Map[int](si, double)

	// Print each value from the filter iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the  iteration was not completed.
	_ = ForEach[int](mi, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
	// 12
	// 14
	// 16
	// 18
	// 20
}

func ExampleReduce() {
	// Average struct used by the average closure.
	type Average struct {
		count   float64
		average float64
	}

	// average reducer closure that is used by Reduce.
	average := func(a Average, v float64) Average {
		a.average = ((a.average * a.count) + v) / (a.count + 1)
		a.count++
		return a
	}

	toFloat := func(i int) float64 {
		return float64(i)
	}

	// Get a sequence iterator that generates values from 1 to 11.
	si := Sequence(1, 11)
	// Map the sequence of ints to float64.
	mi := Map[int, float64](si, toFloat)
	// Reduce all values with reduce closure average to calculate the average. Error is ignored.
	// Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	avg, _ := Reduce[float64](mi, Average{}, average)

	fmt.Println(avg.average)

	// Output:
	// 6
}

func ExampleFromReverseSlice() {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// FromReverseSlice turns the slice into a reverse iterator.
	si := FromReverseSlice(s)

	// Print each value from the slice iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	_ = ForEach[int](si, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 10
	// 9
	// 8
	// 7
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
}

func ExampleStepSequence() {
	// Instead of doing the following we can also use generators.
	// Generators are Iterators that make up data while iterating.
	// This example uses the StepSequence generator.

	// s := []int{1, 3, 5, 7, 9}
	// si := FromSlice(s)

	// Get a sequence iterator that generates values from 1 to 10, but increments with steps of 2.
	si := StepSequence(1, 10, 2)

	// Print each value from the sequence iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	_ = ForEach[int](si, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 1
	// 3
	// 5
	// 7
	// 9
}

func ExampleToSlice() {
	// Iterators can be turned into slices with ToSlice

	// Get a sequence iterator that generates values from 1 to 10, but increments with steps of 2.
	si := StepSequence(1, 10, 2)
	// Convert the iterator into a slice. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	s, _ := ToSlice[int](si)
	// Iterate the slice and print each value.
	for _, v := range s {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 3
	// 5
	// 7
	// 9
}

func ExampleToChannel() {
	// Iterators can send results to a channel with ToChannel.
	c := make(chan int)

	// Get a sequence iterator that generates values from 1 to 10, but increments with steps of 2.
	si := StepSequence(1, 10, 2)
	// Sends the items to a channel. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	go func() {
		defer close(c)
		_ = ToChannel[int](si, c)
	}()
	// Iterate the slice and print each value.
	for v := range c {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 3
	// 5
	// 7
	// 9
}

func ExampleGenerate() {
	// A generator/iterator can be created easily from a closure using Generate.
	// The closure receives two parameters: c is the current count, and r is the amount of items to generate.
	counter := func(p int, c uint64, r uint64) int {
		return p + 1
	}

	gi := Generate[int](0, 3, counter)

	// Print each value from the generating iterator. Error is ignored. Errors can only occur in Iterators which can have
	// an error state. For example a custom iterator that reads data from the database, but the connection is
	// terminated while the iteration was not completed.
	_ = ForEach[int](gi, func(v int) {
		fmt.Println(v)
	})

	// Output:
	// 1
	// 2
	// 3
}

// Tests

type testFixture struct {
	slice                   []int
	resultingIntIterator    Iterable[int]
	resultingStringIterator Iterable[string]
	predicate               PredicateFunc[int]
	mapper                  MapFunc[int, string]
	resultingSlice          []int
	reducer                 ReduceFunc[int, int]
	initialReduceValue      int
	counter                 ForEachFunc[int]
	count                   int
	sum                     int
	generator               GeneratorFunc[string]
	repeat                  uint64
	start                   int
	end                     int
	step                    int
	channel                 chan int
}

var t testFixture

func toSliceOfInts(table *godog.Table) (result []int, err error) {
	var value int
	for _, row := range table.Rows {
		value, err = strconv.Atoi(row.Cells[0].Value)
		if err != nil {
			return
		}
		result = append(result, value)
	}
	return
}

func toSliceOfStrings(table *godog.Table) (result []string) {
	var value string
	for _, row := range table.Rows {
		value = row.Cells[0].Value
		result = append(result, value)
	}
	return
}

func nextReturnsTrueTimesAndThenReturnsFalse(num int) error {
	for ; num > 0; num-- {

		if _, r := t.resultingIntIterator.Next(); r != true {
			return errors.New("expected: true got: false")
		}
	}
	if _, r := t.resultingIntIterator.Next(); r != false {
		return errors.New("expected: false got: true")
	}
	return nil
}

func aSliceIteratorIsReturnedWithIdxContaining(arg1 int) error {
	si := t.resultingIntIterator.(*SliceIterator[int])
	if arg1 != si.idx {
		return fmt.Errorf("expected: %v got: %v", arg1, si.idx)
	}
	return nil
}

func aSliceIteratorIsReturnedWithReverseContainingFalse() error {
	si := t.resultingIntIterator.(*SliceIterator[int])
	if false != si.reverse {
		return fmt.Errorf("expected: %v got: %v", false, si.reverse)
	}
	return nil
}

func aSliceIteratorIsReturnedWithReverseContainingTrue() error {
	si := t.resultingIntIterator.(*SliceIterator[int])
	if true != si.reverse {
		return fmt.Errorf("expected: %v got: %v", true, si.reverse)
	}
	return nil
}

func aSliceIteratorIsReturnedWithValuesContaining(listofints *godog.Table) error {
	si := t.resultingIntIterator.(*SliceIterator[int])
	s, err := toSliceOfInts(listofints)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(s, si.values) {
		return fmt.Errorf("expected: %v got: %v", s, si.values)
	}
	return nil
}

func aSliceWithTheFollowingValues(listofints *godog.Table) (err error) {
	t.slice, err = toSliceOfInts(listofints)
	return
}

func fromSliceIsCalled() {
	t.resultingIntIterator = FromSlice(t.slice)
}

func fromReverseSliceIsCalled() {
	t.resultingIntIterator = FromReverseSlice(t.slice)
}

func aPredicateThatOnlySelectsOddNumbers() {
	t.predicate = func(a int) bool {
		result := (a % 2) != 0
		return result
	}
}

func anIterableWithTheFollowingValues(listofints *godog.Table) error {
	s, err := toSliceOfInts(listofints)
	if err != nil {
		return err
	}
	t.resultingIntIterator = FromSlice(s)
	return nil
}

func filterIsCalled() {
	t.resultingIntIterator = Filter(t.resultingIntIterator, t.predicate)
}

func aMapFunctionThatMultiplesTheValuesAndConvertsTheIntToAStringPrefixedWithTest() {
	t.mapper = func(i int) string {
		return "test" + strconv.Itoa(i*2)
	}
}

func mapIsCalled() {
	t.resultingStringIterator = Map(t.resultingIntIterator, t.mapper)
}

func aSliceIsReturnedWithTheFollowingValues(listofints *godog.Table) error {
	s, err := toSliceOfInts(listofints)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(s, t.resultingSlice) {
		return fmt.Errorf("expected: %v got: %v", s, t.resultingSlice)
	}
	return nil
}

func toSliceIsCalled() (err error) {
	t.resultingSlice, err = ToSlice(t.resultingIntIterator)
	return
}

func aReduceFunctionThatSumsAllValues() {
	t.reducer = func(a, b int) int {
		return a + b
	}
}

func initialValueOf(init int) {
	t.initialReduceValue = init
}

func reduceIsCalled() (err error) {
	t.sum, err = Reduce(t.resultingIntIterator, t.initialReduceValue, t.reducer)
	return
}

func theReturnedSumIs(expected int) error {
	if t.sum != expected {
		return fmt.Errorf("expected: %v got: %v", expected, t.sum)
	}
	return nil
}

func aForeachFunctionThatSumsAndCountsTheCalls() {
	t.counter = func(i int) {
		t.count++
		t.sum += i
	}
}

func foreachIsCalled() error {
	return ForEach(t.resultingIntIterator, t.counter)
}

func theReturnedCountIs(expected int) error {
	if t.count != expected {
		return fmt.Errorf("expected: %v got: %v", expected, t.sum)
	}
	return nil
}

func callingNextUntilFalseIsReturnedShouldReturnTheFollowingStrings(listofints *godog.Table) error {
	expected := toSliceOfStrings(listofints)

	var results []string

	for v, b := t.resultingStringIterator.Next(); b; v, b = t.resultingStringIterator.Next() {
		results = append(results, v)
	}

	if !reflect.DeepEqual(expected, results) {
		return fmt.Errorf("expected: %v got: %v", expected, results)
	}

	return nil
}

func callingNextUntilFalseIsReturnedShouldReturnTheFollowingIntegers(listofints *godog.Table) error {
	expected, err := toSliceOfInts(listofints)

	if err != nil {
		return err
	}

	var results []int

	for v, b := t.resultingIntIterator.Next(); b; v, b = t.resultingIntIterator.Next() {
		results = append(results, v)
	}

	if !reflect.DeepEqual(expected, results) {
		return fmt.Errorf("expected: %v got: %v", expected, results)
	}

	return nil
}

func aGeneratorFuncThatReturnsTheCountAndRepeatConcatenatedWithAComma() {
	t.generator = func(p string, c, r uint64) string {
		return fmt.Sprintf("%d,%d", c, r)
	}
}

func aRepeatValueOf(r int) {
	t.repeat = uint64(r)
}

func generateIsCalled() {
	t.resultingStringIterator = Generate("", t.repeat, t.generator)
}

func aStartValueOf(s int) {
	t.start = s
}

func anEndValueOf(e int) {
	t.end = e
}

func anStepValueOf(s int) {
	t.step = s
}

func stepSequenceIsCalled() {
	t.resultingIntIterator = StepSequence(t.start, t.end, t.step)
}

func sequenceIsCalled() {
	t.resultingIntIterator = Sequence(t.start, t.end)
}

func valuesStringToIntSlice(in string) (result []int, err error) {
	for _, s := range strings.Split(in, ",") {
		i, err2 := strconv.Atoi(s)
		if err2 != nil {
			err = err2
			return
		}
		result = append(result, i)
	}
	return
}

func callingNextUntilFalseIsReturnedShouldReturnTheFollowingValues(values string) error {
	expected, err := valuesStringToIntSlice(values)
	if err != nil {
		return err
	}
	results, err := ToSlice(t.resultingIntIterator)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(expected, results) {
		return fmt.Errorf("expected: %v got: %v", expected, results)
	}
	return nil
}

type ErrorIterator[T any] struct{}

func (e *ErrorIterator[T]) Next() (T, bool) {
	var t T
	return t, false
}

func (e *ErrorIterator[T]) Error() error {
	return errors.New("iterator not implemented")
}

func anIterableInAnErrorState() {
	t.resultingIntIterator = &ErrorIterator[int]{}
}

func errorOfStringIteratorReturnsAnError() error {
	if t.resultingStringIterator.Error() == nil {
		return errors.New("expected an error but got nil")
	}
	return nil
}

func errorOfStringIteratorReturnsNil() error {
	if t.resultingStringIterator.Error() != nil {
		return errors.New("expected nil but got an error")
	}
	return nil
}

func errorOfIntIteratorReturnsAnError() error {
	if t.resultingIntIterator.Error() == nil {
		return errors.New("expected an error but got nil")
	}
	return nil
}

func errorOfIntIteratorReturnsNil() error {
	if t.resultingIntIterator.Error() != nil {
		return errors.New("expected nil but got an error")
	}
	return nil
}

func aClosedChannelWithTheFollowingValues(listofints *godog.Table) {
	t.channel = make(chan int)
	go func() {
		values, err := toSliceOfInts(listofints)
		if err != nil {
			panic(err)
		}
		for _, v := range values {
			t.channel <- v
		}
		close(t.channel)
	}()
}

func fromChannelIsCalled() {
	t.resultingIntIterator = FromChannel(t.channel)
}

func theChannelIsClosed() {
	close(t.channel)
}

func theFollowingValuesAreReceivedOnTheChannel(listofints *godog.Table) error {
	var results []int
	expected, err := toSliceOfInts(listofints)
	if err != nil {
		return err
	}
	for v := range t.channel {
		results = append(results, v)
	}
	if !reflect.DeepEqual(expected, results) {
		return fmt.Errorf("expected: %v got: %v", expected, results)
	}
	return nil
}

func toChannelIsCalled() {
	go func() {
		defer close(t.channel)
		err := ToChannel(t.resultingIntIterator, t.channel)
		if err != nil {
			panic(err)
		}
	}()
}

func aChannel() {
	t.channel = make(chan int)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	t = testFixture{}

	ctx.Step(`^a slice with the following values:$`, aSliceWithTheFollowingValues)
	ctx.Step(`^FromSlice is called$`, fromSliceIsCalled)
	ctx.Step(`^FromReverseSlice is called$`, fromReverseSliceIsCalled)
	ctx.Step(`^Next\(\) returns true (\d+) times and then returns false$`, nextReturnsTrueTimesAndThenReturnsFalse)
	ctx.Step(`^a SliceIterator is returned with \.idx containing (-?\d+)$`, aSliceIteratorIsReturnedWithIdxContaining)
	ctx.Step(`^a SliceIterator is returned with \.reverse containing false$`, aSliceIteratorIsReturnedWithReverseContainingFalse)
	ctx.Step(`^a SliceIterator is returned with \.values containing:$`, aSliceIteratorIsReturnedWithValuesContaining)
	ctx.Step(`^a SliceIterator is returned with \.reverse containing true$`, aSliceIteratorIsReturnedWithReverseContainingTrue)
	ctx.Step(`^a predicate that only selects odd numbers$`, aPredicateThatOnlySelectsOddNumbers)
	ctx.Step(`^an Iterable with the following values:$`, anIterableWithTheFollowingValues)
	ctx.Step(`^Filter is called$`, filterIsCalled)
	ctx.Step(`^a map function that multiples the values and converts the int to a string, prefixed with test$`, aMapFunctionThatMultiplesTheValuesAndConvertsTheIntToAStringPrefixedWithTest)
	ctx.Step(`^Map is called$`, mapIsCalled)
	ctx.Step(`^a slice is returned with the following values:$`, aSliceIsReturnedWithTheFollowingValues)
	ctx.Step(`^ToSlice is called$`, toSliceIsCalled)
	ctx.Step(`^a reduce function that sums all values$`, aReduceFunctionThatSumsAllValues)
	ctx.Step(`^initial value of (\d+)$`, initialValueOf)
	ctx.Step(`^Reduce is called$`, reduceIsCalled)
	ctx.Step(`^Foreach is called$`, foreachIsCalled)
	ctx.Step(`^a foreach function that sums and counts the calls$`, aForeachFunctionThatSumsAndCountsTheCalls)
	ctx.Step(`^The returned count is (\d+)$`, theReturnedCountIs)
	ctx.Step(`^The returned sum is (\d+)$`, theReturnedSumIs)
	ctx.Step(`^calling Next\(\) until false is returned should return the following integers:$`, callingNextUntilFalseIsReturnedShouldReturnTheFollowingIntegers)
	ctx.Step(`^calling Next\(\) until false is returned should return the following strings:$`, callingNextUntilFalseIsReturnedShouldReturnTheFollowingStrings)
	ctx.Step(`^a GeneratorFunc that returns the count and repeat concatenated with a comma\.$`, aGeneratorFuncThatReturnsTheCountAndRepeatConcatenatedWithAComma)
	ctx.Step(`^a repeat value of (\d+)$`, aRepeatValueOf)
	ctx.Step(`^Generate\(\) is called$`, generateIsCalled)
	ctx.Step(`^a start value of (-?\d+)$`, aStartValueOf)
	ctx.Step(`^an end value of (-?\d+)$`, anEndValueOf)
	ctx.Step(`^an step value of (-?\d+)$`, anStepValueOf)
	ctx.Step(`^StepSequence is called$`, stepSequenceIsCalled)
	ctx.Step(`^Sequence is called$`, sequenceIsCalled)
	ctx.Step(`^calling Next\(\) until false is returned should return the following values: "([^"]*)"$`, callingNextUntilFalseIsReturnedShouldReturnTheFollowingValues)
	ctx.Step(`^an Iterable in an error state$`, anIterableInAnErrorState)
	ctx.Step(`^Error\(\) of int iterator returns an error$`, errorOfIntIteratorReturnsAnError)
	ctx.Step(`^Error\(\) of int iterator returns nil$`, errorOfIntIteratorReturnsNil)
	ctx.Step(`^Error\(\) of string iterator returns an error$`, errorOfStringIteratorReturnsAnError)
	ctx.Step(`^Error\(\) of string iterator returns nil$`, errorOfStringIteratorReturnsNil)
	ctx.Step(`^a closed channel with the following values:$`, aClosedChannelWithTheFollowingValues)
	ctx.Step(`^FromChannel is called$`, fromChannelIsCalled)
	ctx.Step(`^the channel is closed$`, theChannelIsClosed)
	ctx.Step(`^the following values are received on the channel$`, theFollowingValuesAreReceivedOnTheChannel)
	ctx.Step(`^ToChannel is called$`, toChannelIsCalled)
	ctx.Step(`^a channel$`, aChannel)

}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// Benchmarks

func BenchmarkFilter(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	odd := func(v int) bool {
		return (v % 2) != 0
	}

	benchFunc := func() []int {
		si := FromSlice(s)
		fi := Filter[int](si, odd)
		ns, _ := ToSlice[int](fi)
		return ns
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}

func BenchmarkFilterInIdiomaticGo(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	benchFunc := func() []int {
		var ns []int
		for _, v := range s {
			if (v % 2) != 0 {
				ns = append(ns, v)
			}
		}
		return ns
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}

func BenchmarkFilterMap(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	odd := func(v int) bool {
		return (v % 2) != 0
	}

	benchFunc := func() []string {
		si := FromSlice(s)
		fi := Filter[int](si, odd)
		mi := Map[int, string](fi, strconv.Itoa)
		ns, _ := ToSlice[string](mi)
		return ns
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}

func BenchmarkFilterMapInIdiomaticGo(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	benchFunc := func() []string {
		var ns []string
		for _, v := range s {
			if (v % 2) != 0 {
				ns = append(ns, strconv.Itoa(v))
			}
		}
		return ns
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}

func BenchmarkFilterMapReduce(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	odd := func(v int) bool {
		return (v % 2) != 0
	}

	join := func(builder *strings.Builder, value string) *strings.Builder {
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(value)
		return builder
	}

	benchFunc := func() string {
		si := FromSlice(s)
		fi := Filter[int](si, odd)
		mi := Map[int, string](fi, strconv.Itoa)
		sb, _ := Reduce[string, *strings.Builder](mi, &strings.Builder{}, join)
		return sb.String()
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}

func BenchmarkFilterMapReduceInIdiomaticGo(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	benchFunc := func() string {
		builder := strings.Builder{}
		for _, v := range s {
			if (v % 2) != 0 {
				if builder.Len() > 0 {
					builder.WriteString(", ")
				}
				builder.WriteString(strconv.Itoa(v))
			}
		}
		return builder.String()
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}

func filterIntSlice(in []int, predicate func(int) bool) (output []int) {
	for _, v := range in {
		if predicate(v) {
			output = append(output, v)
		}
	}
	return
}

func mapIntSliceToStringSlice(in []int, mapper func(int) string) (output []string) {
	for _, v := range in {
		output = append(output, mapper(v))
	}
	return
}

func reduceStringSliceToString(in []string, init *strings.Builder, reducer func(*strings.Builder, string) *strings.Builder) (output string) {
	for _, v := range in {
		init = reducer(init, v)
	}
	return init.String()
}

func BenchmarkFilterMapDIY(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	odd := func(v int) bool {
		return (v % 2) != 0
	}

	join := func(builder *strings.Builder, value string) *strings.Builder {
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(value)
		return builder
	}

	benchFunc := func() string {
		return reduceStringSliceToString(mapIntSliceToStringSlice(filterIntSlice(s, odd), strconv.Itoa), &strings.Builder{}, join)
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}

func BenchmarkFilterMapDIY2(b *testing.B) {

	var s []int

	for n := 0; n < 1000; n++ {
		s = append(s, n)
	}

	odd := func(v int) bool {
		return (v % 2) != 0
	}

	join := func(builder *strings.Builder, value string) *strings.Builder {
		if builder.Len() > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(value)
		return builder
	}

	benchFunc := func() string {
		sb := strings.Builder{}
		for _, v := range s {
			if odd(v) {
				join(&sb, strconv.Itoa(v))
			}
		}
		return sb.String()
	}

	for n := 0; n < b.N; n++ {
		benchFunc()
	}
}
