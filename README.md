# Todo App

Simple todo application built on [Fyne](https://fyne.io/).
This is helpful to me when I need a simple, lightweight application for keeping me organized.
This is pretty much a learning project for me and a way to work on something that solves my immediate problems.
Let me know if it's useful to you too!

If you'd like to report a bug, then you've come to the right place.
Create an issue here and I'll get to it as soon as I have the time.

It supports projects and tasks to help support focus.
* Switch between projects with the dropdown at the top.
* Re-order tasks within a project in a way that makes sense, like to define priority.
* Check tasks as done with the checkbox.
* Update the task summary by double-clicking on the task's text.

## Design priorities
These are the guiding principles with which I'll add features (or not) to the app.
Feel free to create an issue here or make PRs for new features, but they will be accepted or declined based on these principles.

* **Simple**: Tolerate as little complexity as possible, both in the UI and in the code.
* **Extensible**: I would like to be able to add simple extensions to this (see the Plugins section below).
* **Fast**: The application should be highly responsive with absolutely minimal wait times on the part of the user for normal operation.
* **Light**: The application size should get nowhere near a GB. It should be perpetually easy and quick to install and run. 
* **Non-invasive**: The application should require minimal access to the OS and user data without compromising proper operation. Any operation on or collection of user data should be explicitly agreed upon by user consent.

## Data
Data is currently (see SQLite support below) stored as a JSON file in your home directory called `.todo_file`.
This format may change at any time, without notice, so don't try to integrate with this directly.

## Planned changes
There are a few things that I think would be neat to add to the application, that I may or may not do.
Make an issue here for suggestions.

### [Task descriptions](https://github.com/drognisep/todoapp/issues/3)
Add a task description for more extended documentation.
The description should use Fyne's markdown rendering functionality.

### Report a bug in the app
Kind of a no-brainer feature.
This will open an issue in Github with some base boilerplate to report bugs with the app.

This is affected by the **Non-invasive** principle, and explicit consent should be indicated before gathering any environment information.
If the user declines, a bug should still be opened, but without environment details.

### Plugins
I would really like to add the concept of a plugin to the application using [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin), so I can get some experience using the library, and to add some one-off functionality.

This will likely result in adding project level and global level metadata to the data file.

### SQLite support
There are 100% Go ports of the SQLite library now, and I think this would be a better fit for more complex uses of the application.

With this change, I may or may not add a feature to port the JSON model to SQLite.

### A real name
"TODO" is not very descriptive or memorable.
I think this name will change as the application evolves, and I figure out how it's most useful.

### Encrypted at rest
I'm not sure if this is truly necessary, but it would be nice to have.
This would be *much* more necessary to control plugin access to task data, or if there's a plugin authorization model to support the **Non-invasive** principle.
This should require a password to load a key into memory and decrypt the data store.
In a SQLite context, I'm not sure how to accomplish this, so it'll require some research.
