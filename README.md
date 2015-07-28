jsonpatch is a simple little Go library for creating and applying RFC
6902 JSON patches.

Semi-unique among its features is the ability to generate paranoid
patches that include tests that validate that the segments of JSON
being patched have not changed in the time the patch was generated to
the time it was applied.
