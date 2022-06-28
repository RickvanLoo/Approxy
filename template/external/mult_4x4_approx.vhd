-- Salim Ullah, Semeen Rehman, Bharath Srinivas Prabakaran, Florian Kriebel, Muhammad Abdullah Hanif, Muhammad Shafique, and Akash Kumar. 2018. Area-optimized low-latency approximate multipliers for FPGA-based hardware accelerators. In Proceedings of the 55th Annual Design Automation Conference (DAC '18). ACM, New York, NY, USA, Article 159, 6 pages. DOI: https://doi.org/10.1145/3195970.3195996
library IEEE;
use IEEE.STD_LOGIC_1164.ALL;

use IEEE.NUMERIC_STD.ALL;
use IEEE.STD_LOGIC_UNSIGNED.ALL;

library UNISIM;
use UNISIM.VComponents.all;

entity {{.EntityName}} is
generic (word_size: integer:={{.BitSize}}); 
Port ( 
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end {{.EntityName}};

architecture Behavioral of {{.EntityName}} is
signal pp0, pp1 : std_logic_vector(5 downto 0);
signal gen, prop, carries1, output : std_logic_vector(3 downto 0);
signal gen1, prop1, carries2, output1 : std_logic_vector(3 downto 0);
signal prod1 : std_logic_vector(1 downto 0);
begin
lut_inst1: lut6_2 -- P1 and P2
generic map(INIT => X"B4CCF00066AACC00")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => b(0),
I4 => b(1),
I5 => '1',
O5 => pp0(1),
O6 => pp0(2)
);

lut_inst2: lut6_2 -- P3
generic map(INIT => X"C738F0F0FF000000")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => a(3),
I4 => b(0),
I5 => b(1),
O6 => pp0(3)
);

lut_inst3: lut6_2 -- P4
generic map(INIT => X"07C0FF0000000000")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => a(3),
I4 => b(0),
I5 => b(1),
O6 => pp0(4)
);

lut_inst4: lut6_2 -- P5
generic map(INIT => X"F800000000000000")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => a(3),
I4 => b(0),
I5 => b(1),
O6 => pp0(5)
);

pp0(0) <= '0';
pp1(0) <= '0';

lut_inst6: lut6_2 -- P1 and P2
generic map(INIT => X"B4CCF00066AACC00")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => b(2),
I4 => b(3),
I5 => '1',
O5 => pp1(1),
O6 => pp1(2)
);

lut_inst7: lut6_2 -- P3
generic map(INIT => X"C738F0F0FF000000")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => a(3),
I4 => b(2),
I5 => b(3),
O6 => pp1(3)
);



lut_inst81: lut6_2 -- P4
generic map(INIT => X"F800000000000000")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => a(3),
I4 => b(2),
I5 => b(3),
O6 => gen(3)
);


lut_inst8: lut6_2 -- P4
generic map(INIT => X"07C0FF0000000000")
port map(
I0 => a(0),
I1 => a(1),
I2 => a(2),
I3 => a(3),
I4 => b(2),
I5 => b(3),
O6 => prop(3)
);


--------------------
lut_inst10: lut6_2 -- P5
generic map(INIT => X"007F7F80FFEAEA00")
port map(
I0 => pp0(2),
I1 => a(0),
I2 => b(2),
I3 => pp0(3),
I4 => pp1(1),
I5 => '1',
O5 => gen(0),
O6 => prop(0)
);

LUT_inst11: LUT6_2
generic map(INIT => X"6666666688888888")
port map(
	I0 => pp0(4),
	I1 => pp1(2),
	I2 => '1',
	I3 => '1',
	I4 => '1',
	I5 => '1',
	O5 => gen(1),
	O6 => prop(1)
	);
	
LUT_inst12: LUT6_2
generic map(INIT => X"6666666688888888")
port map(
	I0 => pp0(5),
	I1 => pp1(3),
	I2 => '1',
	I3 => '1',
	I4 => '1',
	I5 => '1',
	O5 => gen(2),
	O6 => prop(2)
	);


carry_inst1: CARRY4
port map (
	DI => gen,
	S => prop,
	O => output,
	CO => carries1,
	CI => '0',
	CYINIT => '0'
	);

prod(6 downto 3) <= output(3 downto 0);
prod(7) <= carries1(3);
prod(1) <= pp0(1);

-------------------------------------
lut_inst14: lut6_2 -- P1 and P2
generic map(INIT => X"5FA05FA088888888")
port map(
I0 => a(0),
I1 => b(0),
I2 => b(2),
I3 => pp0(2),
I4 => '1',
I5 => '1',
O5 => prod(0),
O6 => prod(2)
);

end Behavioral;
