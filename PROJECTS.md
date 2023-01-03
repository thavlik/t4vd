## Static Kata Recognition
Task: given a *single frame*, determine the particular technique being employed.

## Kata Stream Analysis
Using the static kata recognition model to convert a video's frames into an event stream, analyze the sequences of kata across all videos to yield a directed graph representing the statistical likelihood of transition from any given kata to another.

This problem is computationally similar to [sequence homology](https://en.wikipedia.org/wiki/Sequence_homology) in genetics. Instead of analyzing DNA to model the likelihood of a given nucleotide being the next in sequence, the stream of kata from the video is analyzed to determine the likelihood of a given kata being next.

## Dynamic Kata Recognition
Task: given a *short video*, determine the particular technique being employed.

This project is exponentially more difficult than static recognition, as determining when exactly a kata video starts and ends is nontrivial.

## Kata Identification & Naming
Task: analyze neural networks trained on kata recognition to discover kata that exist in the wild but have escaped naming.

Despite pioneering the technique, [Jason Von Flue](https://en.wikipedia.org/wiki/Jason_Von_Flue) was most likely not the first human to ever perform the [Von Flue Choke](https://www.youtube.com/watch?v=rkwpb7RBu90). Similarly, there are other kata that are only rarely performed, have not yet been popularized, and thus are unnamed.

This project hypothesizes that innumerable kata exist that have not yet been named, for a variety of reasons:
- Naming is not helpful
- The technique is not effective
- The kata is a blend of other kata

## Submission Sequence Identification & Naming
At a higher level of analysis, individual techniques are strung together into a single, cohesive technique in what may be called a *submission sequence*. As opposed to detecting these sequences on video directly, which may transpire over the course of several minutes, they can be found by identifying loops in a directed graph formed by analyzing sequences of kata across all videos.

