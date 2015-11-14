# FsWatcher : A golang file changes tracker
FsWatcher as its name states (or not) is a simple file system change watcher.

FsWatcher does not rely on any optimized native system artifacts to check changes but just polls the monitored directories for file changes.

Please note that this behavior is not adapted for really huge directory trees and is better suited at project/workspaces sized trees.

FsWatcher is easy to use (please see godoc for more details):

	cwd, err := os.Getwd()
	w, e := fswatcher.NewFsWatcher(cwd, 500)
	if e != nil {
		...
	}

	w.Skip(".git")
	w.RegisterFileExtension(".go", func(path string, info os.FileInfo) bool {
		...
		return false
	})
	w.Watch()
    
*FsWatcher was initially developped for the [gobot](http://github.com/fjecker/gobot) project and was written while learning/discovering go* 
