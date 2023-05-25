# Workers as Actors

A simple example. Demo steps:

1. Run the base workflow, start it multiple times. Hooray, you have loosely coupled computation units.

2. Add signals. Now we can receive messages.
2a. As a result, start a child workflow, then send it a message. Hooray, sending too!

3. Kill all the workers. Restart them. Things continue.

4. Continue as new. 