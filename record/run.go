package record

// Run drives the emulator step-by-step for the given number of steps without
// wall-clock timing. stepFn executes one instruction. renderFn, if non-nil,
// updates the framebuffer pixel data from display memory and is called before
// video capture at declared steps. video may be nil.
func Run(
	stepFn func(),
	rec *Recorder,
	video *VideoRecorder,
	steps int,
	renderFn func(),
) {
	for range steps {
		rec.Step(stepFn)

		if video != nil {
			step := rec.CurrentStep()
			if video.NeedsCapture(step) {
				if renderFn != nil {
					renderFn()
				}
				video.Observe(step)
			}
		}
	}
}
