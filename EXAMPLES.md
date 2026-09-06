# Unit-Aware Arithmetic - Examples

Real-world examples demonstrating the library's capabilities.

## Table of Contents
- [Physics Examples](#physics-examples)
- [Engineering Examples](#engineering-examples)
- [Scientific Computing](#scientific-computing)
- [Everyday Calculations](#everyday-calculations)
- [Error Prevention](#error-prevention)

---

## Physics Examples

### Classical Mechanics

```python
from dimensional import Quantity, units

# Calculate velocity
distance = Quantity(100, units.meter)
time = Quantity(9.58, units.second)
velocity = distance / time
print(f"Usain Bolt's speed: {velocity}")  # 10.438 m/s

# Calculate acceleration
initial_velocity = Quantity(0, units.meter_per_second)
final_velocity = Quantity(30, units.meter_per_second)
time = Quantity(5, units.second)
acceleration = (final_velocity - initial_velocity) / time
print(f"Acceleration: {acceleration}")  # 6.0 m/s²

# Calculate force (F = ma)
mass = Quantity(1500, units.kilogram)  # Car mass
acceleration = Quantity(2.5, units.meter_per_second_squared)
force = mass * acceleration
print(f"Force required: {force}")  # 3750 N
```

### Energy and Power

```python
from dimensional import Quantity, units

# Kinetic energy (KE = 1/2 * m * v²)
mass = Quantity(1000, units.kilogram)
velocity = Quantity(25, units.meter_per_second)  # 90 km/h
kinetic_energy = 0.5 * mass * (velocity ** 2)
print(f"Kinetic energy: {kinetic_energy}")  # 312500 J

# Power consumption
energy = Quantity(3600, units.joule)
time = Quantity(1, units.hour)
power = energy / time
print(f"Power: {power}")  # 1.0 W

# Convert to kilojoules
energy_kj = energy.to(units.kilojoule)
print(f"Energy: {energy_kj}")  # 3.6 kJ
```

### Gravitational Calculations

```python
from dimensional import Quantity, units

# Free fall
g = Quantity(9.8, units.meter_per_second_squared)
time = Quantity(3, units.second)
distance = 0.5 * g * (time ** 2)
print(f"Distance fallen: {distance}")  # 44.1 m

# Potential energy (PE = mgh)
mass = Quantity(50, units.kilogram)
height = Quantity(10, units.meter)
potential_energy = mass * g * height
print(f"Potential energy: {potential_energy}")  # 4900 J
```

---

## Engineering Examples

### Electrical Engineering

```python
from dimensional import Quantity, units

# Power (P = F·v)
force = Quantity(100, units.newton)
velocity = Quantity(10, units.meter_per_second)
power = force * velocity
print(f"Power: {power}")  # 1000 W

# Ohm's Law (V = IR) - future feature
# voltage = current * resistance
```

### Mechanical Engineering

```python
from dimensional import Quantity, units

# Pressure calculation
force = Quantity(1000, units.newton)
area = Quantity(0.5, units.meter) * Quantity(0.5, units.meter)
pressure = force / area
print(f"Pressure: {pressure}")  # 4000 Pa

# Torque calculation (τ = F × r)
force = Quantity(100, units.newton)
radius = Quantity(0.5, units.meter)
torque = force * radius
print(f"Torque: {torque}")  # 50 J (dimensionally equivalent to N·m)
```

### Fluid Dynamics

```python
from dimensional import Quantity, units

# Flow rate (Q = A × v)
area = Quantity(0.01, units.meter) * Quantity(0.01, units.meter)
velocity = Quantity(2, units.meter_per_second)
flow_rate = area * velocity
print(f"Flow rate: {flow_rate}")  # 0.0002 m³/s

# Reynolds number (dimensionless) - future feature
# Re = (ρ × v × L) / μ
```

---

## Scientific Computing

### Chemistry

```python
from dimensional import Quantity, units

# Temperature conversion
temp_celsius = Quantity(25, units.celsius)
temp_kelvin = temp_celsius.to(units.kelvin)
print(f"Room temperature: {temp_kelvin}")  # 298.15 K

# Concentration (mol/m³) - future feature
```

### Astronomy

```python
from dimensional import Quantity, units

# Distance a light signal travels in 4.24 days
speed_of_light = Quantity(299792458, units.meter_per_second)
time = Quantity(4.24, units.day)
distance = speed_of_light * time
distance_km = distance.to(units.kilometer)
print(f"Distance: {distance_km}")  # ~1.10e11 km
```

---

## Everyday Calculations

### Cooking

```python
from dimensional import Quantity, units

# Recipe scaling
original_flour = Quantity(500, units.gram)
scale_factor = 1.5
scaled_flour = original_flour * scale_factor
print(f"Flour needed: {scaled_flour}")  # 750 g

# Temperature conversion
oven_fahrenheit = Quantity(350, units.fahrenheit)
oven_celsius = oven_fahrenheit.to(units.celsius)
print(f"Oven temperature: {oven_celsius}")  # 176.67 °C
```

### Travel

```python
from dimensional import Quantity, units

# Speed conversion (mph to km/h)
speed = Quantity(65, units.mile) / Quantity(1, units.hour)
speed_kmh = speed.to(units.kilometer_per_hour)
print(f"Speed limit: {speed_kmh}")  # 104.6 km/h

# Distance conversion
distance_miles = Quantity(500, units.mile)
distance_km = distance_miles.to(units.kilometer)
print(f"Trip distance: {distance_km}")  # 804.67 km
```

### Fitness

```python
from dimensional import Quantity, units

# Running pace
distance = Quantity(5, units.kilometer)
time = Quantity(25, units.minute)
pace = time / distance
print(f"Pace: {pace}")  # 5 min/km

# Energy burn rate
energy = Quantity(500, units.kilojoule)
time = Quantity(1, units.hour)
burn_rate = energy / time
print(f"Burn rate: {burn_rate}")  # ~138.9 W
```

---

## Error Prevention

### Catching Common Mistakes

```python
from dimensional import Quantity, units, IncompatibleUnitsError

# Example 1: Adding incompatible units
try:
    distance = Quantity(100, units.meter)
    time = Quantity(10, units.second)
    result = distance + time  # ERROR!
except IncompatibleUnitsError as e:
    print(f"Error caught: {e}")
    # Output: Cannot add m and s: incompatible dimensions

# Example 2: Comparing incompatible units
try:
    mass = Quantity(50, units.kilogram)
    length = Quantity(2, units.meter)
    if mass > length:  # ERROR!
        pass
except IncompatibleUnitsError as e:
    print(f"Error caught: {e}")
    # Output: Cannot compare kg and m

# Example 3: Wrong unit conversion
try:
    temperature = Quantity(25, units.celsius)
    distance = temperature.to(units.meter)  # ERROR!
except IncompatibleUnitsError as e:
    print(f"Error caught: {e}")
    # Output: Cannot convert °C to m: incompatible dimensions
```

### Correct Usage

```python
from dimensional import Quantity, units

# Correct: Adding compatible units
d1 = Quantity(100, units.meter)
d2 = Quantity(50, units.centimeter)
total = d1 + d2  # Automatic conversion!
print(total)  # 100.5 m

# Correct: Comparing compatible units
mass1 = Quantity(1, units.kilogram)
mass2 = Quantity(1000, units.gram)
print(mass1 == mass2)  # True (automatic conversion)

# Correct: Unit conversion
temp_f = Quantity(32, units.fahrenheit)
temp_c = temp_f.to(units.celsius)
print(temp_c)  # 0.0 °C
```

---

## Advanced Examples

### Dimensional Analysis

```python
from dimensional import Quantity, units

# Derive units through calculation
mass = Quantity(10, units.kilogram)
velocity = Quantity(5, units.meter_per_second)
time = Quantity(2, units.second)

# Force = mass × acceleration
# acceleration = velocity / time
acceleration = velocity / time
force = mass * acceleration

print(f"Acceleration: {acceleration}")  # 2.5 m/s²
print(f"Force: {force}")  # 25 N

# Verify dimensions
print(f"Force dimension: L={force.unit.dimension.length}, "
      f"M={force.unit.dimension.mass}, "
      f"T={force.unit.dimension.time}")
# Output: L=1, M=1, T=-2 (correct for force!)
```

### Unit Consistency Checking

```python
from dimensional import Quantity, units

def calculate_work(force: Quantity, distance: Quantity) -> Quantity:
    """Calculate work done (W = F × d)."""
    # The library automatically checks that force has correct dimensions
    work = force * distance
    
    # Verify result has energy dimensions (L²MT⁻²)
    assert work.unit.dimension.length == 2
    assert work.unit.dimension.mass == 1
    assert work.unit.dimension.time == -2
    
    return work

# Usage
force = Quantity(50, units.newton)
distance = Quantity(10, units.meter)
work = calculate_work(force, distance)
print(f"Work done: {work}")  # 500 J
```

---

## Best Practices

### 1. Always Use Units

```python
# ❌ Bad: Raw numbers
distance = 100  # What unit?
time = 10       # What unit?
speed = distance / time  # 10 what?

# ✅ Good: Explicit units
distance = Quantity(100, units.meter)
time = Quantity(10, units.second)
speed = distance / time  # 10 m/s (clear!)
```

### 2. Convert Early

```python
# ❌ Bad: Convert at the end
def calculate_speed(distance_miles, time_hours):
    speed = distance_miles / time_hours
    return speed * 1.60934  # Convert to km/h

# ✅ Good: Use units throughout
def calculate_speed(distance: Quantity, time: Quantity) -> Quantity:
    return distance / time
```

### 3. Validate Inputs

```python
def calculate_pressure(force: Quantity, area: Quantity) -> Quantity:
    """Calculate pressure (P = F/A)."""
    # Validate dimensions
    if force.unit.dimension.mass != 1 or force.unit.dimension.length != 1:
        raise ValueError("Force must have dimensions of mass × length / time²")
    
    if area.unit.dimension.length != 2:
        raise ValueError("Area must have dimensions of length²")
    
    return force / area
```

---

**Last Updated**: 2026-08-28
**More Examples**: See test files for comprehensive examples
