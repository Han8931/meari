package curriculum

// Pre-authored Physics curriculum. Physics exercises are prose/math responses:
// the executor verifies that a response was submitted, while the lesson content
// and tutor feedback guide the learner's reasoning.

func physicsBeginner() []Module {
	return []Module{
		{
			Name: "Motion and Forces",
			Topics: []Topic{
				{
					ID:    "phys-b-motion",
					Title: "Motion, speed and velocity",
					Lesson: `Motion means position changes over time. Speed tells how fast distance is
covered; velocity includes direction. Walking 3 meters east and then 3 meters
west gives 6 meters of distance but 0 meters of displacement because you end
where you started.

Analogy: distance is the full path of a maze; displacement is the straight
arrow from start to finish. A car speedometer shows speed, but a navigation app
uses velocity because direction matters.

Example:
    A runner travels 100 m in 20 s.
    average speed = distance / time = 100 / 20 = 5 m/s

If the runner returns to the start after that trip, total distance is still
100 m, but displacement is 0 m. Physics keeps those ideas separate because they
answer different questions: "how much road?" versus "where did you end up?"`,
					Challenge: Challenge{
						Prompt:      "A student walks 8 m east, then 3 m west in 5 s. Find total distance, displacement, average speed, and average velocity. Explain the difference in words.",
						StarterCode: "Write your reasoning here:\n\n",
						Solution:    "Distance is 8 + 3 = 11 m. Displacement is 5 m east because the final position is 5 m east of the start. Average speed is 11/5 = 2.2 m/s. Average velocity is 5/5 = 1 m/s east. Distance counts the whole path; displacement is the net change in position.",
					},
				},
				{
					ID:    "phys-b-newton",
					Title: "Newton's laws",
					Lesson: `Newton's laws describe how forces change motion.

1. Inertia: objects keep their motion unless a net force acts. A book stays on
a table; a puck slides until friction slows it.
2. F = ma: acceleration depends on net force and mass. The same push gives a
shopping cart more acceleration than a loaded truck.
3. Action-reaction: forces come in pairs on different objects. Your foot pushes
the ground backward; the ground pushes you forward.

Analogy: mass is like "stubbornness" against changes in motion. A bowling ball
is more stubborn than a tennis ball. Force is the attempt to change motion;
acceleration is the result.

Example:
    net force = 12 N, mass = 3 kg
    acceleration = F/m = 12/3 = 4 m/s^2`,
					Challenge: Challenge{
						Prompt:      "A 2 kg cart experiences 10 N right and 4 N left. Find the net force and acceleration. Name which Newton law you used.",
						StarterCode: "Write your reasoning here:\n\n",
						Solution:    "The forces oppose, so net force is 10 - 4 = 6 N to the right. Using Newton's second law F = ma, acceleration is a = F/m = 6/2 = 3 m/s^2 to the right.",
					},
				},
				{
					ID:    "phys-b-energy",
					Title: "Energy and conservation",
					Lesson: `Energy is the ability to cause change. Kinetic energy is energy of motion;
gravitational potential energy is stored because of height.

    kinetic energy: KE = 1/2 mv^2
    gravitational potential: PE = mgh

Conservation of energy says energy changes form but does not disappear. A ball
held high has gravitational potential energy. As it falls, that energy becomes
kinetic energy. With friction or air resistance, some mechanical energy becomes
thermal energy, like hands warming when rubbed together.

Analogy: energy is like money moving between accounts. It can move from
"height account" to "motion account" to "heat account", but the bookkeeping
must still balance.`,
					Challenge: Challenge{
						Prompt:      "Explain what happens to the energy of a skateboarder rolling down a ramp. Include potential energy, kinetic energy, and friction.",
						StarterCode: "Write your explanation here:\n\n",
						Solution:    "At the top, the skateboarder has more gravitational potential energy. As they roll down, height decreases and speed increases, so potential energy changes into kinetic energy. Friction and air resistance transform some energy into thermal energy and sound, so not all potential energy becomes useful motion.",
					},
				},
				{
					ID:    "phys-b-waves",
					Title: "Waves, sound and light",
					Lesson: `A wave carries energy without carrying matter all the way with it. In water,
the surface moves up and down while the wave travels sideways.

Important wave ideas:
    amplitude: size of the disturbance
    wavelength: distance from crest to crest
    frequency: cycles per second
    speed = frequency x wavelength

Sound is a pressure wave in matter; it needs a medium like air or water. Light
is an electromagnetic wave and can travel through empty space.

Analogy: a stadium wave travels around the stadium, but each person mostly just
stands up and sits down. The pattern moves; the people do not run around the
stadium.`,
					Challenge: Challenge{
						Prompt:      "A wave has frequency 5 Hz and wavelength 2 m. Find its speed. Then explain why sound cannot travel through empty space but light can.",
						StarterCode: "Write your reasoning here:\n\n",
						Solution:    "Wave speed is f times wavelength: 5 x 2 = 10 m/s. Sound needs particles to compress and spread out, so it cannot travel in empty space. Light is an electromagnetic wave, so it can travel through vacuum.",
					},
				},
			},
		},
	}
}

func physicsIntermediate() []Module {
	return []Module{
		{
			Name: "Classical Physics",
			Topics: []Topic{
				{
					ID:    "phys-i-kinematics",
					Title: "Kinematics with graphs",
					Lesson: `Kinematics describes motion without asking what caused it. Graphs are a
powerful language for motion.

On a position-time graph, slope is velocity. A steep line means fast motion. A
horizontal line means no change in position. On a velocity-time graph, slope is
acceleration and area is displacement.

Analogy: a graph is a motion story compressed into a picture. Position-time
answers "where are you?" Velocity-time answers "how fast and which way?"

For constant acceleration:
    v = v0 + at
    x = x0 + v0 t + 1/2 a t^2
    v^2 = v0^2 + 2a(x - x0)

Example: dropping from rest with g = 9.8 m/s^2 for 2 s gives
v = 19.6 m/s downward and distance = 1/2 gt^2 = 19.6 m.`,
					Challenge: Challenge{
						Prompt:      "A ball starts from rest and accelerates at 3 m/s^2 for 4 s. Find final velocity and displacement. Explain what the area under the velocity-time graph represents.",
						StarterCode: "Write your reasoning here:\n\n",
						Solution:    "Final velocity is v = 0 + 3*4 = 12 m/s. Displacement is x = 0 + 0 + 1/2*3*4^2 = 24 m. On a velocity-time graph, area represents displacement; here the area is a triangle with base 4 and height 12, so 24 m.",
					},
				},
				{
					ID:    "phys-i-momentum",
					Title: "Momentum and collisions",
					Lesson: `Momentum is mass times velocity:
    p = mv

Momentum is conserved in an isolated system. That means total momentum before
an interaction equals total momentum after. This is why recoil happens: if a
skater throws a heavy ball forward, the skater rolls backward so total momentum
balances.

Analogy: momentum is motion bookkeeping. During a collision, objects can trade
momentum like players passing a ball, but the team's total stays fixed if no
outside force acts.

Kinetic energy may or may not be conserved. In elastic collisions it is; in
inelastic collisions some kinetic energy becomes heat, sound, deformation, or
internal energy.`,
					Challenge: Challenge{
						Prompt:      "A 2 kg cart moving 3 m/s sticks to a 1 kg cart at rest. Find their shared final speed using momentum conservation.",
						StarterCode: "Write your reasoning here:\n\n",
						Solution:    "Initial momentum is 2*3 + 1*0 = 6 kg m/s. After sticking, total mass is 3 kg. Conservation gives 6 = 3v, so v = 2 m/s in the original direction.",
					},
				},
				{
					ID:    "phys-i-electricity",
					Title: "Circuits and fields",
					Lesson: `Electric charge creates electric fields, and fields exert forces on charges.
In circuits, voltage is energy per charge, current is flow of charge, and
resistance opposes current.

Ohm's law:
    V = IR

Analogy: a simple circuit is like water in pipes. Voltage is like pressure,
current is like flow rate, and resistance is like a narrow pipe. The analogy is
not perfect, but it helps: more pressure gives more flow; more resistance gives
less flow.

Power is energy per time:
    P = IV = I^2 R = V^2/R

Series resistors add directly. Parallel paths give current more routes, so
equivalent resistance decreases.`,
					Challenge: Challenge{
						Prompt:      "A 9 V battery is connected to a 3 ohm resistor. Find the current and power. Explain using the water-flow analogy.",
						StarterCode: "Write your reasoning here:\n\n",
						Solution:    "Using V = IR, current is I = V/R = 9/3 = 3 A. Power is P = IV = 3*9 = 27 W. In the water analogy, the 9 V battery is pressure and the resistor is a narrow pipe; lower resistance allows more charge flow.",
					},
				},
				{
					ID:    "phys-i-thermo",
					Title: "Temperature, heat and entropy",
					Lesson: `Temperature measures average microscopic kinetic energy. Heat is energy
transferred because of temperature difference. They are related but not the
same: a spark has high temperature but little total heat; a warm bathtub has
lower temperature but much more thermal energy.

The first law of thermodynamics is energy conservation:
    change in internal energy = heat added - work done by the system

Entropy measures how spread out energy is and how many microscopic arrangements
fit the same large-scale state. Natural processes tend to spread energy out.

Analogy: a tidy room has few arrangements; a messy room has many. Entropy is
not simply "mess", but the room analogy helps explain why mixed states are more
likely than perfectly organized states.`,
					Challenge: Challenge{
						Prompt:      "Explain why an ice cube melts in warm water using heat transfer and entropy. Avoid saying 'cold moves'.",
						StarterCode: "Write your explanation here:\n\n",
						Solution:    "Thermal energy transfers from warmer water to colder ice because of the temperature difference. The ice's internal energy increases until its structure breaks into liquid water. The combined system moves toward a state where energy is more spread out, which corresponds to higher entropy.",
					},
				},
			},
		},
	}
}

func physicsAdvanced() []Module {
	return []Module{
		{
			Name: "Modern and University Physics",
			Topics: []Topic{
				{
					ID:    "phys-a-rotational",
					Title: "Rotation and torque",
					Lesson: `Rotational motion mirrors linear motion:
    position -> angle
    velocity -> angular velocity
    mass -> moment of inertia
    force -> torque

Torque measures how effectively a force causes rotation:
    tau = r F sin(theta)

Analogy: opening a door is easier when you push far from the hinge. The same
force gives more torque because the lever arm is longer. A figure skater spins
faster when pulling arms inward because moment of inertia decreases and angular
momentum is conserved.

Rotational kinetic energy is:
    KE_rot = 1/2 I omega^2`,
					Challenge: Challenge{
						Prompt:      "Why does pushing near a door handle work better than pushing near the hinge? Explain using torque and lever arm.",
						StarterCode: "Write your explanation here:\n\n",
						Solution:    "Torque depends on force times lever arm. Near the handle, the distance from the hinge is large, so the same force produces more torque. Near the hinge, the lever arm is small, so the door gets much less rotational effect.",
					},
				},
				{
					ID:    "phys-a-quantum",
					Title: "Quantum ideas",
					Lesson: `Quantum physics describes systems where energy, measurement, and probability
cannot be treated like everyday objects. Light and matter show both wave-like
and particle-like behavior.

A photon has energy:
    E = hf

The wavefunction is not a little physical wave in space like water. It is a
probability amplitude: when squared, it gives probabilities for measurement
outcomes.

Analogy: before rolling a die, probability describes possible outcomes, but the
quantum version is deeper because amplitudes can interfere. Two paths can
combine to make an outcome more likely or cancel to make it less likely, like
overlapping ripples in water.`,
					Challenge: Challenge{
						Prompt:      "Explain the double-slit experiment in terms of probability amplitudes. Why is it misleading to say an electron is just a tiny marble?",
						StarterCode: "Write your explanation here:\n\n",
						Solution:    "In the double-slit experiment, amplitudes from different paths combine and interfere, creating a pattern of likely and unlikely detection spots. A tiny marble would go through one slit with a definite path and would not produce the same interference pattern. The electron is detected in localized events, but its probabilities evolve wave-like.",
					},
				},
				{
					ID:    "phys-a-relativity",
					Title: "Relativity and spacetime",
					Lesson: `Special relativity starts from two ideas: the laws of physics are the same in
all inertial frames, and the speed of light in vacuum is the same for all
observers. The result is not that "everything is relative"; instead, space and
time adjust so light speed stays invariant.

Consequences:
    moving clocks run slow relative to an observer
    moving lengths contract along the motion direction
    simultaneity depends on reference frame

Analogy: north and east are different directions, but rotating a map mixes them.
Space and time are different, but changing reference frame mixes measurements
of space and time into spacetime.`,
					Challenge: Challenge{
						Prompt:      "Use the train-and-lightning idea to explain why simultaneity can depend on the observer's frame of reference.",
						StarterCode: "Write your explanation here:\n\n",
						Solution:    "If lightning strikes both ends of a train, an observer on the ground midway between strike points may receive both flashes together and call them simultaneous. An observer on the moving train is moving toward one flash and away from the other, while light speed is the same in both directions, so they receive the flashes at different times and do not call the strikes simultaneous.",
					},
				},
				{
					ID:    "phys-a-fields",
					Title: "Fields and Maxwell's equations",
					Lesson: `A field assigns a physical quantity to every point in space and time. A
temperature map assigns temperature; an electric field assigns force per charge.
Fields let physics describe influence locally instead of imagining mysterious
instant action at a distance.

Maxwell's equations describe electric and magnetic fields. In words:
    electric charges create electric fields
    there are no isolated magnetic charges in ordinary electromagnetism
    changing magnetic fields create electric fields
    currents and changing electric fields create magnetic fields

Together they predict electromagnetic waves: changing electric and magnetic
fields sustain each other and travel at the speed of light.`,
					Challenge: Challenge{
						Prompt:      "Explain how a changing magnetic field can create an electric field, and why that idea matters for generators or transformers.",
						StarterCode: "Write your explanation here:\n\n",
						Solution:    "A changing magnetic field induces an electric field, which can push charges in a wire and create current. In a generator, motion changes magnetic flux through coils, producing electric power. In a transformer, changing current in one coil creates changing magnetic field, which induces voltage in another coil.",
					},
				},
				{
					ID:    "phys-a-lagrangian",
					Title: "Lagrangian mechanics",
					Lesson: `Newtonian mechanics focuses on forces. Lagrangian mechanics focuses on energy
and paths. The Lagrangian is usually:
    L = kinetic energy - potential energy

The actual path makes the action stationary, where action is the time integral
of L. This sounds abstract, but it is powerful because it works beautifully with
constraints, generalized coordinates, fields, and advanced physics.

Analogy: imagine comparing every possible route a marble could take. Nature's
route is not chosen by a tiny planner; the mathematics says neighboring
possible paths balance so the action is stationary. Like a soap film finding a
minimal surface, the final behavior emerges from a global condition.`,
					Challenge: Challenge{
						Prompt:      "Compare Newton's F = ma approach with the Lagrangian approach for a pendulum. Why can energy-based coordinates be easier?",
						StarterCode: "Write your explanation here:\n\n",
						Solution:    "Newton's approach tracks forces and components, including tension, which can be awkward for circular motion. The Lagrangian approach can use the pendulum angle as the coordinate and write kinetic and potential energy directly. The constraint from the string is built into the coordinate choice, so the math focuses on the actual degree of freedom.",
					},
				},
			},
		},
	}
}
