---
name: Debug a graphics problem
description: Enable an agent to debug graphics problems
---

# Debug a graphics problem

If you need to see what might be wrong with graphics, you can look at the
debug file that ends in `.screen`. This file is the screen log. It records
frames with a full copy of the graphics on screen in text form. These frames
are sampled roughly once per second.

Each dot on the screen is represented by a character in the frame.

When debugging anything, several other files are helpful:

- the metrics file (how many times certain events occur; ends in `.metrics`)
- the instruction log (all instructions sorted by place in memory it was
  executed; file ends in `.asm`)
- the timeset file (how much time each instruction took; ends in `.time`)
