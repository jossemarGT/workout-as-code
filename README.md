# Workout as code

> Are you a tired developer with a bunch of workouts in many places? Get WAC!

Workout as code (WAC) is a simple but structured way of writing workouts as
plain text. No more bookmarks on your social media apps :muscle:

## Getting started

The `wac` cli is under incubation, but the files it uses are human-readable and
ready to go. You can find the workouts under the `data/` directory.


Pull requests are welcome, you can find the [contribution guidelines here](docs/CONTRIBUTING.md).

### Using the workout V0 language

Inspired by the
[workout of the day whiteboards](https://www.google.com/search?q=wod+whiteboard&tbm=isch)
you can find on any cross/functional training centre. The workout language aims
to describe a WOD in a simple way.

Keep in mind that a "workout" is a set of exercises, each one having a specific
amount of repetitions or target time. A minimal workout would look like this.

```markdown
1 Push up
```

As shown above, you have the number of repetitions followed by the exercise
name. The amount can be immidiately followed by a suffix, i.e.: `1m`, but the
are no semantic rules over these, so it is up to the human that reads it to
interpret it.

Although it is grammatically correct (in workout language terms), a single
exercise with one repetition does not help anybody. So let's improve our minimal
workout.

```markdown
10  Push-ups
1m  Rest
10  Push-ups
1m  Rest
8   Australian pull-ups
1m  Rest
8   Australian pull-ups
1m  Rest
```

That is better, but we might have a problem with the this example, as we would
prefer to group the push exercises apart from the pull ones. In that case, we
can use an H2 title to delimit each group of exercises. These groups are called
"circuits". Our minimal workout now should look like this


```markdown
## Push

10  Push-ups
1m  Rest
10  Push-ups
1m  Rest

## Pull

8   Australian pull-ups
1m  Rest
8   Australian pull-ups
1m  Rest
```

Woth mentioning, a circuit can hold metadata, like limits or additional rules;
this can be done by adding it after the ` - ` keyword (note the empty spaces
around the dash). A good example of circuit metadata is this Crossfit staple
workout:

```markdown
## Mary - AMRAP 20m

5	Handstand push-ups
10	Pistol squats, alterning
15	Pull-ups
```

At the moment, `wac` cannot understand circuit metadata semantics; still, this
can be useful for other humans ;).

Finally, it is not required but an excellent gesture to add a title to the whole
workout. This one can hold metadata as well, in case it is needed.

```markdown
# // The sample workout //

10m mobility warm up

## Push set

10	Push ups
1m	Rest
10	Push ups
1m	Rest
10	Push ups
1m	Rest


## // Pull set //

8   Austrlian pull ups
1m  Rest
8   Austrlian pull ups
1m  Rest
8   Austrlian pull ups
1m  Rest

## Core set - 4x rest: 1m

10x	Crunches
1m	Rest

## Cardio set - 4x

10x Jumping jacks
1m  Rest
```

### The workout v0 specification

This languange aims to be a Markdown super set, the current ebnf represantion
goes like this:

```ebnf
  Workout = (Set | Circuit)* .
  Set = <quantity> <gstring> .
  Circuit = Title+ Set* .
  Title = <titlestart> TitleFragment* (<metadiv> Metadata+)? <titleend> .
  TitleFragment = (<tgstring> | <puntc>) .
  Metadata = (<wodtype> | <quantity> | MetaWord) .
  MetaWord = <ident> (<colon> MetaTag)? .
  MetaTag = (<quantity> | <ident>) .
```

In case you would like to deep dive on this topic please visit the `pkg`
directory.

## v0 Roadmap

- [x] Write down the idea (more or less)
- [x] Dump initial workouts as plain files. **Human oriented** and machine readable
- [x] Prototype wac tool- Lint workout files
- [ ] Prototype wac tool - Export workouts as "Seconds" mobile app input file
- [ ] *Maybe* separate the tool src code from workouts data
- [ ] ???
- [ ] Profit!
