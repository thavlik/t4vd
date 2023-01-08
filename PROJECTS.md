# List of Brazilian Jiu-jitsu Machine Learning Projects

## Static Kata Recognition
**Task:** given a *single frame*, determine the particular technique being employed.

**Status:** Active

The current plan to solve this task is to use the app to present the user with a frame and have them tag it with one or more key words. No major roadblocks are foreseen in solving this task, and it is currently the highest priority.

## Kata Stream Analysis
**Task:** using the static kata recognition model to convert a video's frames into an event stream, analyze the sequences of kata across all videos to yield a directed graph representing the statistical likelihood of transition from any given kata to another.

**Status:** awaiting progress on *Static Kata Recognition*

This problem is computationally similar to [sequence homology](https://en.wikipedia.org/wiki/Sequence_homology) in genetics. Instead of analyzing DNA to model the likelihood of a given nucleotide being the next in sequence, the stream of kata from the video is analyzed to determine the likelihood of a given kata being next.

## Dynamic Kata Recognition
**Task:** given a *short video*, determine the particular technique being employed.

**Status:** backburner

This project is exponentially more difficult than static recognition, as determining when exactly a kata video starts and ends is nontrivial.

## Kata Identification & Naming
**Task:** analyze the neural networks trained in *Static Kata Recognition* to discover kata that exist in the wild but have escaped naming.

**Status:** backburner

Despite pioneering the technique, [Jason Von Flue](https://en.wikipedia.org/wiki/Jason_Von_Flue) was most likely not the first human to ever perform the [Von Flue Choke](https://www.youtube.com/watch?v=rkwpb7RBu90). Similarly, there are other kata that are only rarely performed, have not yet been popularized, and thus are unnamed.

This project hypothesizes that innumerable kata exist that have not yet been named, for a variety of reasons:
- Naming is not helpful
- The technique is not effective
- The kata is a blend of other kata

## Submission Sequence Identification & Naming
**Task:** identify and name submission sequences as loops in directed graphs yielded by the *Kata Stream Analysis* project.

**Status:** backburner

At a higher level of analysis, individual techniques are strung together into a single, cohesive technique in what may be called a *submission sequence*. As opposed to detecting these sequences on video directly, which may transpire over the course of several minutes, they can be found by identifying loops in a directed graph formed by analyzing sequences of kata across all videos.

Many of these loops should already been well known, such as *armbar-omoplata-triangle*.

## Simple Belt Detection (Gi only)
**Task**: given a single frame of two opponents in gi, determine their belt colors.

**Status:** backburner

## Landmark Detection
**Task:** predict locations and orientation of opponents' joints a la [On-Device, Real-Time Hand Tracking with MediaPipe
](https://ai.googleblog.com/2019/08/on-device-real-time-hand-tracking-with.html).

**Status:** backburner

This is believed to be an incredibly difficult task, but that solving it precipitate great progress on a variety of other tasks.

## Area of Interest Detection
**Task:** crop a frame to include only the grapplers

**Status:** backburner

There are plans to implement a cropping tool in the app to facilitate data creation for this task. Solving it should increase data efficiency on other tasks.

## Per-Athlete Kata Sequence Analysis & Comparison
**Task:** compare a set/distribution of kata sequences to another set/distribution of kata sequences 

**Status:** backburner

The edges on the kata sequence graph represent likelihoods of observing the transition. Transitioning from one kata to another can then be modeled as a univariate gaussian, which can be generalized to multivariate gaussian as one kata can transition into many. The average overlap between an individual's transition distributions and the average grappler's can be computed as a proxy of how "unique" a grappler's techniques are.

The overall complexity of an athlete's kata transition graph (as represented by the number of edges, loops, etc.) can then computed to objectively quantify how "complex" a grappler's techniques are.

The ability to compare an individual's most recent kata sequence graph to one in their past may prove a useful training tool.

### Personalized Kata Graph Art
**Task:** create a force-directed graph from an individual's kata sequence graph. Integrate said force directed graph into beautiful art that evolves over time in response to the changes observed in their kata sequences.

**Status:** conceptual phase

## Vector-embedded Kata
**Task:** create a semantically rich, low-dimensional vector embedding for all kata.

**Status:** conceptual phase

As word tags can be applied to kata to assign them discrete categories, so too can they be unsupervisedly clustered. The resultant latent space can then be called "kata space", or **K**-space. A good K-space should have large areas dedicated to named kata.

Analysis of K-space embeddings should illuminate variations in technique that make or break it.

A specific instance of kata can then be analyzed to compare joint positions/orientations against an established model. 

Features of certain kata, such as the switch between threading & choking arms in D'arce, should be appreciable in K-space.

## Searchable Kata Video Database
**Task:** create a searchable database of technique videos, with unrestricted input terms (e.g. "closed guard into arm bar into scissor sweep")

**Status:** conceptual phase

## Automatic Keyframing
**Task:** train a model to automatically impose the grapplers over a green screen background

**Status:** conceptual phase

The creation of supervised data for this task would most likely mirror that of [neurosurgery-video-dataset](https://github.com/thavlik/neurosurgery-video-dataset)'s.

## Grappler Identity Detection
**Task:** given a short video of two grapplers at the start of their roll, paint each grapplers' pixels with a unique color to separately identify them. 

**Status:** conceptual phase

Like *Automatic Keyframing*, the creation of supervised data for this task would most likely mirror that of [neurosurgery-video-dataset](https://github.com/thavlik/neurosurgery-video-dataset)'s. Creating labels for this task would thus create labels for both.

## Dual-track Kata Sequence Streams
**Task:** disambiguate a single stream of kata sequences for both grapplers into two grappler-specific kata sequence streams

**Status:** conceptual phase

Solving this task would allow for distinctions to be made between e.g. side control top vs. bottom, attacking the armbar vs. defending it. Kata sequence streams of this type are expected to be semantically richer. 

Progress on this task would most likely require *Grappler Identity Detection* to be solved.

## Animated 3D Models
**Task:** use generative modeling to animate 3D characters performing techniques

**Status:** conceptual phase

An obvious improvement to a kata sequence visualization tool would be the inclusion of videos demonstrating the technique whenever it appears in the graph. As opposed to searching a database of existing videos, novel videos could be synthesized using a variety of methods. One such method utilizes skeleton and triangle mesh intermediates. The skeleton may be yielded by solving *Landmark Detection*.

## BJJ Chess
**Task:** derive a decision tree from kata sequence analysis and create a turn-based video game where the efficacy of each action is computed from the action's likelihood of success in the real world

**Status:** conceptual phase

## Biomechanical Analysis
**Task:** given a rich representation of limb motion/orientation, estimate the force vectors utilized by a particular technique.

**Status:** conceptual phase

Possible solutions include generative modeling where the force vectors predict the next frame.

