package fto

var workingSetTemplate = []Session{
	{
		Set{
			Percent: 65,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 75,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 85,
			Reps:    5,
			AMRAP:   true,
			Type:    Working,
		},
	},
	{
		Set{
			Percent: 70,
			Reps:    3,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 80,
			Reps:    3,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 90,
			Reps:    3,
			AMRAP:   true,
			Type:    Working,
		},
	},
	{
		Set{
			Percent: 75,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 85,
			Reps:    3,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 95,
			Reps:    1,
			AMRAP:   true,
			Type:    Working,
		},
	},
}

var deloadTemplate = map[DeloadType]Session{
	Deload1: {
		Set{
			Percent: 40,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 50,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 60,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
	},
	Deload2: {
		Set{
			Percent: 50,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 60,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 70,
			Reps:    5,
			AMRAP:   false,
			Type:    Working,
		},
	},
	Deload3: {
		Set{
			Percent: 65,
			Reps:    3,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 75,
			Reps:    3,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 85,
			Reps:    3,
			AMRAP:   false,
			Type:    Working,
		},
	},
	Deload4: {
		Set{
			Percent: 40,
			Reps:    10,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 50,
			Reps:    8,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 60,
			Reps:    6,
			AMRAP:   false,
			Type:    Working,
		},
	},
	Deload5: {
		Set{
			Percent: 50,
			Reps:    10,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 60,
			Reps:    8,
			AMRAP:   false,
			Type:    Working,
		},
		Set{
			Percent: 70,
			Reps:    6,
			AMRAP:   false,
			Type:    Working,
		},
	},
}
