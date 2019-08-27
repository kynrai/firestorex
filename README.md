# FirestoreX

Extra firestore helper functions that wrap existing firestore calls in an idiomatic way to provide conveniance functions when working with Google firestore.

## Auto concurrent batches

Write a large struct into firestore without worrying about about the 500 document limit, also run concurrent batches to dramatically speed up writes.
