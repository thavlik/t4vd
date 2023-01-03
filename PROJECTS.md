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
