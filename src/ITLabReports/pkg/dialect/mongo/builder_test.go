package mongo

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFunc_Builder_Predicate(t *testing.T) {
	t.Run("Or_And_Equal", func(t *testing.T) {
		{
			expected := bson.M{
				"$or": bson.A{
					bson.M{
						"field_1": 12,
					},
					bson.M{
						"$and": bson.A{
							bson.M{
								"field_3": 15,
							},
							bson.M{
								"field_4": 17,
							},
						},
					},
				},
			}
			actual := Or(
				EQ("field_1", 12),
				And(
					EQ("field_3", 15),
					EQ("field_4", 17),
				),
			).raw

			require.Equal(t, expected, actual)
		}
	})

	t.Run(
		"Like",
		func(t *testing.T) {
			actual := Like("name", "some", I).raw
			expected := bson.M{
				"name": bson.M{
					"$regex": "some",
					"$options": "i",
				},
			}

			require.Equal(
				t,
				expected,
				actual,
			)
		},
	)

	t.Run("not", func(t *testing.T) {
		expected := bson.M{
			"$not": bson.M{
				"field_1": 12,
			},
		}
		actual := P().Not(P().EQ("field_1", 12)).raw

		require.Equal(t, expected, actual)
	})

	t.Run("And_NotEqual", func(t *testing.T) {
		expected := bson.M{
			"$and": bson.A{
				bson.M{
					"field": bson.M{
						"$ne": 1,
					},
				},
				bson.M{
					"field": bson.M{
						"$ne": 12,
					},
				},
			},
		}

		actual := And(
			NEQ("field", 1),
			NEQ("field", 12),
		).raw

		require.Equal(t, expected, actual)
	})

	t.Run("In", func(t *testing.T) {
		expected := bson.M{
			"field": bson.M{
				"$in": bson.A{
					1, 2, 3,
				},
			},
		}

		actual := In("field", 1, 2, 3).raw
		require.Equal(t, expected, actual)
	})

	t.Run("NotIn", func(t *testing.T) {
		expected := bson.M{
			"field": bson.M{
				"$nin": bson.A{
					1, 2, 3,
				},
			},
		}

		actual := NotIn("field", 1, 2, 3).raw
		require.Equal(t, expected, actual)
	})

	t.Run("Field_LT_GTE", func(t *testing.T) {
		expected := bson.M{
			"field": bson.M{
				"$lt": 3,
				"$gt": 5,
			},
		}

		actual := P().LT("field", 3).GT("field", 5).raw
		require.Equal(t, expected, actual)
	})

	t.Run("ElemMatch", func(t *testing.T) {
		expect := bson.M{
			"arrayField": bson.M{
				"$elemMatch": bson.M{
					"field_1": 2,
					"$or": bson.A{
						bson.M{
							"field_2": bson.M{
								"$gte": 12,
							},
						},
						bson.M{
							"field_2": bson.M{
								"$lt": -10,
							},
						},
					},
				},
			},
		}

		actual := ElemMatch(
			"arrayField",
			Or(
				GTE("field_2", 12),
				LT("field_2", -10),
			).EQ("field_1", 2),
		).raw

		require.Equal(t, expect, actual)
	})

	t.Run("ID_IN", func(t *testing.T) {
		{
			expect := bson.M{
				"_id": bson.M{
					"$in": bson.A{
						1, 2, 3, 4,
					},
				},
			}
	
			actual := P().In("_id", 1, 2).In("_id", 3, 4).BSON()
	
			require.Equal(t, expect, actual)
		}
		{
			expect := bson.M{
				"_id": bson.M{
					"$in": bson.A{
						1, 2, 3, 4,
					},
				},
			}

			ids := []interface{}{1, 2, 3, 4}
			actual := P().In("_id", ids...).BSON()
			require.Equal(t, expect, actual)
		}
	})

	t.Run(
		"And_Multiply",
		func(t *testing.T) {
			exprected := bson.M{
				"$and": bson.A{
					bson.M{
						"name": "some",
					},
					bson.M{
						"name": "some_2",
					},
				},
			}

			actual := P().And(P().EQ("name", "some"), P().EQ("name", "some_2")).Object()

			require.Equal(
				t,
				exprected,
				actual,
			)
		},
	)

	t.Run(
		"Or_Multiply",
		func(t *testing.T) {
			exprected := bson.M{
				"$or": bson.A{
					bson.M{
						"name": "some",
					},
					bson.M{
						"name": "some_2",
					},
				},
			}

			actual := P().Or(P().EQ("name", "some"), P().EQ("name", "some_2")).Object()

			require.Equal(
				t,
				exprected,
				actual,
			)
		},
	)
}

func TestFunc_SetBuilder(t *testing.T) {
	expected := bson.M{
		"$set": bson.M{
			"some_field": 1,
			"some_object_field": bson.M{
				"field_1": 14,
				"field_2": "str",
			},
			"some_slice_field": bson.A{
				bson.M{
					"field_1": 7,
					"field_2": 12,
				},
				"string",
			},
			"sliceField.$":                10,
			"anotherSliceField.$.field_1": "updated_field",
		},
	}

	actual := Set().SetField("some_field", 1).
		SetField(
			"some_object_field",
			Object().AddField("field_1", 14).AddField("field_2", "str").Object(),
		).
		SetField(
			"some_slice_field",
			Array().
				AddElem(
					Object().
						AddField("field_1", 7).
						AddField("field_2", 12).
						Object(),
				).
				AddElem("string").
				Array(),
		).
		SetSliceValue("sliceField", 10).
		SetSliceField("anotherSliceField", "field_1", "updated_field").raw

	require.Equal(t, expected, actual)
}

func TestFunc_Update_Set(t *testing.T) {
	expectFilter := bson.M{
		"key_1": 12,
		"key_2": 13,
	}

	expectUpdate := bson.M{
		"$set": bson.M{
			"field_1": 13,
		},
	}

	update := Update("some_collection").
				Filter(EQ("key_1", 12).EQ("key_2", 13)).
				Set(
					func(s *SetBuidler) {
						s.SetField("field_1", 13)
					},
				)
	
	require.Equal(t, expectFilter, update.filter.raw)
	require.Equal(t, expectUpdate, update.update)
}

func TestFunc_Pipeline_Builder(t *testing.T) {
	expectPipeline := []bson.M{
		{
			"$match": bson.M{
				"field_1": 12,
				"field_2": "value_2",
			},
		},
		{
			"$projects": bson.M{
				"field_1": 0,
				"field_2": 1,
			},
		},
	}

	pipeline := Pipeline().AddStage(
		func(o *ObjectBuilder) {
			o.AddField(
				"$match",
				Object().
					AddField(
						"field_1",
						12,
					).
					AddField(
						"field_2",
						"value_2",
					),
			)
		},
	).AddStage(
		func(o *ObjectBuilder) {
			o.AddField(
				"$projects",
				Object().
					AddField("field_1", 0).
					AddField("field_2", 1),
			)
		},
	).ToPipeline()

	require.Equal(t, expectPipeline, pipeline)
}

func TestFunc_Map(t *testing.T) {
	expected := bson.M{
		"$map": bson.M{
			"input": bson.M{
				"$range": bson.A{
					0,
					10,
				},
			},
			"in": bson.M{
				"$substrCP": bson.A{
					"$name",
					"$$this",
					1,
				},
			},
		},
	}

	actual := Map(
		MapBody{
			Input: Range(
				0,
				10,
				nil,
			),
			In: Object().AddSliceField("$substrCP", "$name", "$$this", 1),
		},
	).Object()

	require.Equal(t, expected, actual)
}

func TestFunc_AddFieldsStage(t *testing.T) {
	expect := []bson.M{
		{
			"$addFields": bson.M{
				"field_1": "str",
				"field_2": bson.M{
					"field_1": 12,
				},
				"field_3": bson.A{
					1,
					2,
					3,
				},
			},
		},
	}

	actual := Pipeline().AddFields(
		func(o *ObjectBuilder) {
			o.AddField("field_1", "str")
		},
		func(o *ObjectBuilder) {
			o.AddField("field_2", Object().AddField("field_1", 12))
		},
		func(o *ObjectBuilder) {
			o.AddField("field_3", Array().AddElems(1,2,3))
		},
	).ToPipeline()

	require.Equal(t, expect, actual)
}

func TestFunc_Filter(t *testing.T) {
	expct := bson.M{
		"$filter": bson.M{
			"input": "$field_array",
			"as": "num",
			"cond": bson.M{
				"$and": bson.A{
					bson.M{
						"$gte": bson.A{
							"field_1", 12,
						},
					},
					bson.M{
						"$lt": bson.A{
							"field_2", 10,
						},
					},
				},
			},
		},
	}


	actual := Filter(
		FilterBody{
			Input: "$field_array",
			As: "num",
			Cond: And(
				AGTE(
					"field_1",
					12,
				),
				ALT(
					"field_2",
					10,
				),
			),
		},
	).Object()

	require.Equal(t, expct, actual)
}

func TestFunc_ObjectBuilder_AddFieldIf(t *testing.T) {
	expect := bson.M{
		"field_1": bson.M{
			"field_2": 13,
		},
	}

	actual := Object().AddField("field_1", Object().AddFieldIf("field_2", 13, func() bool {return true}).AddFieldIf("field_3", 14, func() bool {return false})).Object()

	require.Equal(t, expect, actual)
}

func TestFunc_Func(t *testing.T) {
	
}