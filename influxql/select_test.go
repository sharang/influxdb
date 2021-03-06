package influxql_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/influxdata/influxdb/influxql"
	"github.com/influxdata/influxdb/pkg/deep"
)

// Second represents a helper for type converting durations.
const Second = int64(time.Second)

// Ensure a SELECT min() query can be executed.
func TestSelect_Min(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		if !reflect.DeepEqual(opt.Expr, MustParseExpr(`min(value)`)) {
			t.Fatalf("unexpected expr: %s", spew.Sdump(opt.Expr))
		}

		return influxql.NewCallIterator(&FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},
		}}, opt), nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT min(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT distinct() query can be executed.
func TestSelect_Distinct_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 1 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 11 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 12 * Second, Value: 2},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT distinct(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 20}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 1 * Second, Value: 19}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 5 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT distinct() query can be executed.
func TestSelect_Distinct_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 1 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 11 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 12 * Second, Value: 2},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT distinct(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 20}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 1 * Second, Value: 19}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 5 * Second, Value: 10}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT distinct() query can be executed.
func TestSelect_Distinct_String(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &StringIterator{Points: []influxql.StringPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: "a"},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 1 * Second, Value: "b"},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: "c"},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: "b"},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: "d"},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 11 * Second, Value: "d"},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 12 * Second, Value: "d"},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT distinct(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: "a"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 1 * Second, Value: "b"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 5 * Second, Value: "c"}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: "d"}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT mean() query can be executed.
func TestSelect_Mean_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT mean(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 3.2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT mean() query can be executed.
func TestSelect_Mean_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT mean(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 3.2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT median() query can be executed.
func TestSelect_Median_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT median(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 3}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT median() query can be executed.
func TestSelect_Median_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT median(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2.5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 3}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT top() query can be executed.
func TestSelect_Top_NoTags_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT top(value, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 20}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 5}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 4}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT top() query can be executed.
func TestSelect_Top_NoTags_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT top(value, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 20}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 5}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 4}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT top() query can be executed with tags.
func TestSelect_Top_Tags_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT top(value, host, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 10, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "B"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Time: 30 * Second, Value: 100, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "A"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Time: 30 * Second, Value: 5, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "B"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT top() query can be executed with tags.
func TestSelect_Top_Tags_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT top(value, host, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 10, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "B"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 30 * Second, Value: 100, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "A"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 30 * Second, Value: 5, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "B"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT top() query can be executed with tags and group by.
func TestSelect_Top_GroupByTags_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT top(value, host, 1) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY region, time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: 19, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: 100, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: "A"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT top() query can be executed with tags and group by.
func TestSelect_Top_GroupByTags_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT top(value, host, 1) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY region, time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: 19, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: 100, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: "A"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT bottom() query can be executed.
func TestSelect_Bottom_NoTags_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT bottom(value, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 3}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT bottom() query can be executed.
func TestSelect_Bottom_NoTags_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT bottom(value, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 2}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 3}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 1}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 30 * Second, Value: 2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT bottom() query can be executed with tags.
func TestSelect_Bottom_Tags_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT bottom(value, host, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 2, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 10, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "B"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Time: 30 * Second, Value: 1, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "B"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Time: 30 * Second, Value: 100, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "A"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT bottom() query can be executed with tags.
func TestSelect_Bottom_Tags_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT bottom(value, host, 2) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 2, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 10, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 0 * Second, Value: "B"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 30 * Second, Value: 1, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "B"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Time: 30 * Second, Value: 100, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Time: 30 * Second, Value: "A"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT bottom() query can be executed with tags and group by.
func TestSelect_Bottom_GroupByTags_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT bottom(value, host, 1) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY region, time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: 2, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: 3, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: 1, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: "B"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT bottom() query can be executed with tags and group by.
func TestSelect_Bottom_GroupByTags_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3, Aux: []interface{}{"A"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100, Aux: []interface{}{"A"}},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4, Aux: []interface{}{"B"}},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5, Aux: []interface{}{"B"}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT bottom(value, host, 1) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY region, time(30s) fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: 2, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=east"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: 3, Aux: []interface{}{"A"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 0 * Second, Value: "A"},
		},
		{
			&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: 1, Aux: []interface{}{"B"}},
			&influxql.StringPoint{Name: "cpu", Tags: ParseTags("region=west"), Time: 30 * Second, Value: "B"},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT query with a fill(null) statement can be executed.
func TestSelect_Fill_Null_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12 * Second, Value: 2},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT mean(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-01T00:01:00Z' GROUP BY host, time(10s) fill(null)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 40 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 50 * Second, Nil: true}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT query with a fill(<number>) statement can be executed.
func TestSelect_Fill_Number_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12 * Second, Value: 2},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT mean(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-01T00:01:00Z' GROUP BY host, time(10s) fill(1)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20 * Second, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 40 * Second, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 50 * Second, Value: 1}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT query with a fill(previous) statement can be executed.
func TestSelect_Fill_Previous_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("host=A"), Time: 12 * Second, Value: 2},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT mean(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-01T00:01:00Z' GROUP BY host, time(10s) fill(previous)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 20 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 40 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 50 * Second, Value: 2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT stddev() query can be executed.
func TestSelect_Stddev_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT stddev(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 0.7071067811865476}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 0.7071067811865476}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 1.5811388300841898}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT stddev() query can be executed.
func TestSelect_Stddev_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT stddev(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 0.7071067811865476}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 0.7071067811865476}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Nil: true}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 1.5811388300841898}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT stddev() query can be executed.
func TestSelect_Stddev_String(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &StringIterator{Points: []influxql.StringPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: "a"},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: "b"},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: "c"},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: "d"},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: "e"},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: "f"},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: "g"},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: "h"},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: "i"},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: "j"},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: "k"},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT stddev(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: ""}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: ""}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: ""}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: ""}},
		{&influxql.StringPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: ""}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT spread() query can be executed.
func TestSelect_Spread_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT spread(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 0}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 1}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 0}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 4}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT spread() query can be executed.
func TestSelect_Spread_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 1},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 5},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT spread(value) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 1}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 0}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 1}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 0}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 4}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT percentile() query can be executed.
func TestSelect_Percentile_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 9},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 8},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 7},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 54 * Second, Value: 6},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 55 * Second, Value: 5},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 56 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 57 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 58 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 59 * Second, Value: 1},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT percentile(value, 90) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 20}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 3}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 9}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT percentile() query can be executed.
func TestSelect_Percentile_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},

			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 50 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 51 * Second, Value: 9},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 52 * Second, Value: 8},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 53 * Second, Value: 7},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 54 * Second, Value: 6},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 55 * Second, Value: 5},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 56 * Second, Value: 4},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 57 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 58 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 59 * Second, Value: 1},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT percentile(value, 90) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 20}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 3}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
		{&influxql.IntegerPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 50 * Second, Value: 9}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a simple raw SELECT statement can be executed.
func TestSelect_Raw(t *testing.T) {
	// Mock two iterators -- one for each value in the query.
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		if !reflect.DeepEqual(opt.Aux, []string{"v1", "v2"}) {
			t.Fatalf("unexpected options: %s", spew.Sdump(opt.Expr))

		}
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Time: 0, Aux: []interface{}{float64(1), nil}},
			{Time: 1, Aux: []interface{}{nil, float64(2)}},
			{Time: 5, Aux: []interface{}{float64(3), float64(4)}},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT v1, v2 FROM cpu`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{
			&influxql.FloatPoint{Time: 0, Value: 1},
			&influxql.FloatPoint{Time: 0, Nil: true},
		},
		{
			&influxql.FloatPoint{Time: 1, Nil: true},
			&influxql.FloatPoint{Time: 1, Value: 2},
		},
		{
			&influxql.FloatPoint{Time: 5, Value: 3},
			&influxql.FloatPoint{Time: 5, Value: 4},
		},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

// Ensure a SELECT binary expr queries can be executed as floats.
func TestSelect_BinaryExpr_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		makeAuxFields := func(value float64) []interface{} {
			aux := make([]interface{}, len(opt.Aux))
			for i := range aux {
				aux[i] = value
			}
			return aux
		}
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Time: 0 * Second, Value: 20, Aux: makeAuxFields(20)},
			{Name: "cpu", Time: 5 * Second, Value: 10, Aux: makeAuxFields(10)},
			{Name: "cpu", Time: 9 * Second, Value: 19, Aux: makeAuxFields(19)},
		}}, nil
	}

	for _, test := range []struct {
		Name      string
		Statement string
		Points    [][]influxql.Point
	}{
		{
			Name:      "rhs binary add",
			Statement: `SELECT value + 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 22}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 12}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 21}},
			},
		},
		{
			Name:      "lhs binary add",
			Statement: `SELECT 2 + value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 22}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 12}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 21}},
			},
		},
		{
			Name:      "two variable binary add",
			Statement: `SELECT value + value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 40}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 20}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 38}},
			},
		},
		{
			Name:      "rhs binary multiply",
			Statement: `SELECT value * 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 40}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 20}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 38}},
			},
		},
		{
			Name:      "lhs binary multiply",
			Statement: `SELECT 2 * value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 40}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 20}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 38}},
			},
		},
		{
			Name:      "two variable binary multiply",
			Statement: `SELECT value * value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 400}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 100}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 361}},
			},
		},
		{
			Name:      "rhs binary subtract",
			Statement: `SELECT value - 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 18}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 8}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 17}},
			},
		},
		{
			Name:      "lhs binary subtract",
			Statement: `SELECT 2 - value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: -18}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: -8}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: -17}},
			},
		},
		{
			Name:      "two variable binary subtract",
			Statement: `SELECT value - value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 0}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 0}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 0}},
			},
		},
		{
			Name:      "rhs binary division",
			Statement: `SELECT value / 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 10}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 5}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: float64(19) / 2}},
			},
		},
		{
			Name:      "lhs binary division",
			Statement: `SELECT 38 / value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 1.9}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 3.8}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 2}},
			},
		},
		{
			Name:      "two variable binary division",
			Statement: `SELECT value / value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.FloatPoint{Name: "cpu", Time: 0 * Second, Value: 1}},
				{&influxql.FloatPoint{Name: "cpu", Time: 5 * Second, Value: 1}},
				{&influxql.FloatPoint{Name: "cpu", Time: 9 * Second, Value: 1}},
			},
		},
	} {
		itrs, err := influxql.Select(MustParseSelectStatement(test.Statement), &ic, nil)
		if err != nil {
			t.Errorf("%s: parse error: %s", test.Name, err)
		} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, test.Points) {
			t.Errorf("%s: unexpected points: %s", test.Name, spew.Sdump(a))
		}
	}
}

// Ensure a SELECT binary expr queries can be executed as integers.
func TestSelect_BinaryExpr_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		makeAuxFields := func(value int64) []interface{} {
			aux := make([]interface{}, len(opt.Aux))
			for i := range aux {
				aux[i] = value
			}
			return aux
		}
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Time: 0 * Second, Value: 20, Aux: makeAuxFields(20)},
			{Name: "cpu", Time: 5 * Second, Value: 10, Aux: makeAuxFields(10)},
			{Name: "cpu", Time: 9 * Second, Value: 19, Aux: makeAuxFields(19)},
		}}, nil
	}

	for _, test := range []struct {
		Name      string
		Statement string
		Points    [][]influxql.Point
	}{
		{
			Name:      "rhs binary add",
			Statement: `SELECT value + 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 22}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 12}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 21}},
			},
		},
		{
			Name:      "lhs binary add",
			Statement: `SELECT 2 + value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 22}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 12}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 21}},
			},
		},
		{
			Name:      "two variable binary add",
			Statement: `SELECT value + value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 40}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 20}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 38}},
			},
		},
		{
			Name:      "rhs binary multiply",
			Statement: `SELECT value * 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 40}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 20}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 38}},
			},
		},
		{
			Name:      "lhs binary multiply",
			Statement: `SELECT 2 * value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 40}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 20}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 38}},
			},
		},
		{
			Name:      "two variable binary multiply",
			Statement: `SELECT value * value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 400}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 100}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 361}},
			},
		},
		{
			Name:      "rhs binary subtract",
			Statement: `SELECT value - 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 18}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 8}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 17}},
			},
		},
		{
			Name:      "lhs binary subtract",
			Statement: `SELECT 2 - value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: -18}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: -8}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: -17}},
			},
		},
		{
			Name:      "two variable binary subtract",
			Statement: `SELECT value - value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 0}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 0}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 0}},
			},
		},
		{
			Name:      "rhs binary division",
			Statement: `SELECT value / 2 FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 10}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 5}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 9}},
			},
		},
		{
			Name:      "lhs binary division",
			Statement: `SELECT 38 / value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 1}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 3}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 2}},
			},
		},
		{
			Name:      "two variable binary division",
			Statement: `SELECT value / value FROM cpu`,
			Points: [][]influxql.Point{
				{&influxql.IntegerPoint{Name: "cpu", Time: 0 * Second, Value: 1}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 5 * Second, Value: 1}},
				{&influxql.IntegerPoint{Name: "cpu", Time: 9 * Second, Value: 1}},
			},
		},
	} {
		itrs, err := influxql.Select(MustParseSelectStatement(test.Statement), &ic, nil)
		if err != nil {
			t.Errorf("%s: parse error: %s", test.Name, err)
		} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, test.Points) {
			t.Errorf("%s: unexpected points: %s", test.Name, spew.Sdump(a))
		}
	}
}

// Ensure a SELECT (...) query can be executed.
func TestSelect_ParenExpr(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		if !reflect.DeepEqual(opt.Expr, MustParseExpr(`min(value)`)) {
			t.Fatalf("unexpected expr: %s", spew.Sdump(opt.Expr))
		}

		return influxql.NewCallIterator(&FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 11 * Second, Value: 3},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 31 * Second, Value: 100},
		}}, opt), nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT (min(value)) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 19}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 0 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 30 * Second, Value: 100}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}

	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 0 * Second, Value: 20},
			{Name: "cpu", Tags: ParseTags("region=west,host=A"), Time: 1 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=west,host=B"), Time: 5 * Second, Value: 10},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 9 * Second, Value: 19},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 10 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 11 * Second, Value: 2},
			{Name: "cpu", Tags: ParseTags("region=east,host=A"), Time: 12 * Second, Value: 2},
		}}, nil
	}

	// Execute selection.
	itrs, err = influxql.Select(MustParseSelectStatement(`SELECT (distinct(value)) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-02T00:00:00Z' GROUP BY time(10s), host fill(none)`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 0 * Second, Value: 20}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 1 * Second, Value: 19}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=B"), Time: 5 * Second, Value: 10}},
		{&influxql.FloatPoint{Name: "cpu", Tags: ParseTags("host=A"), Time: 10 * Second, Value: 2}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestSelect_Derivative_Float(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &FloatIterator{Points: []influxql.FloatPoint{
			{Name: "cpu", Time: 0 * Second, Value: 20},
			{Name: "cpu", Time: 4 * Second, Value: 10},
			{Name: "cpu", Time: 8 * Second, Value: 19},
			{Name: "cpu", Time: 12 * Second, Value: 3},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT derivative(value, 1s) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-01T00:00:16Z'`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Time: 4 * Second, Value: -2.5}},
		{&influxql.FloatPoint{Name: "cpu", Time: 8 * Second, Value: 2.25}},
		{&influxql.FloatPoint{Name: "cpu", Time: 12 * Second, Value: -4}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}

func TestSelect_Derivative_Integer(t *testing.T) {
	var ic IteratorCreator
	ic.CreateIteratorFn = func(opt influxql.IteratorOptions) (influxql.Iterator, error) {
		return &IntegerIterator{Points: []influxql.IntegerPoint{
			{Name: "cpu", Time: 0 * Second, Value: 20},
			{Name: "cpu", Time: 4 * Second, Value: 10},
			{Name: "cpu", Time: 8 * Second, Value: 19},
			{Name: "cpu", Time: 12 * Second, Value: 3},
		}}, nil
	}

	// Execute selection.
	itrs, err := influxql.Select(MustParseSelectStatement(`SELECT derivative(value, 1s) FROM cpu WHERE time >= '1970-01-01T00:00:00Z' AND time < '1970-01-01T00:00:16Z'`), &ic, nil)
	if err != nil {
		t.Fatal(err)
	} else if a := Iterators(itrs).ReadAll(); !deep.Equal(a, [][]influxql.Point{
		{&influxql.FloatPoint{Name: "cpu", Time: 4 * Second, Value: -2.5}},
		{&influxql.FloatPoint{Name: "cpu", Time: 8 * Second, Value: 2.25}},
		{&influxql.FloatPoint{Name: "cpu", Time: 12 * Second, Value: -4}},
	}) {
		t.Fatalf("unexpected points: %s", spew.Sdump(a))
	}
}
