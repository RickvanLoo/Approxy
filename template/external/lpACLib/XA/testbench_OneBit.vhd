-- "lpACLib" is a library for Low-Power Approximate Computing Modules.
-- Copyright (C) 2016, Walaa El-Harouni, Muhammad Shafique, CES, KIT.
—- E-mail: walaa.elharouny@gmail.com, swahlah@yahoo.com

-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.

-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
-- GNU General Public License for more details.

-- You should have received a copy of the GNU General Public License
-- along with this program. If not, see <http://www.gnu.org/licenses/>.

--------------------------------------------
-- Testbench for all 1-bit FA 
-- Author: Walaa El-Harouni  
--------------------------------------------
LIBRARY ieee;
USE ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

library VITAL;
use VITAL.all;

ENTITY testbench IS
END testbench;

ARCHITECTURE behavior OF testbench IS
	
	COMPONENT AdderAccurateOneBit
	port (A   	 : in std_logic;
	      B 	 : in std_logic;
	      Cin 	 : in std_logic;
	      Sum 	 : out std_logic;
	      Cout   : out std_logic );
	end component;

	signal a: std_logic;
	signal b: std_logic;
	signal cin: std_logic;
	signal sum: std_logic;
	signal cout: std_logic;

BEGIN
	uut : AdderAccurateOneBit port map (A => a, B => b, Cin => cin, Sum => sum, Cout => cout);
		
	stim_proc : process
	begin
		for l in 0 to 95000 loop
		wait for 2 ns;
		a <= '0';
		b <= '0';
		cin <= '0';

		wait for 3 ns;
		a <= '0';
		b <= '0';
		cin <= '1';

		wait for 3 ns;
		a <= '0';
		b <= '1';
		cin <= '0';

		wait for 3 ns;
		a <= '0';
		b <= '1';
		cin <= '1';

		wait for 3 ns;
		a <= '1';
		b <= '0';
		cin <= '0';

		wait for 3 ns;
		a <= '1';
		b <= '0';
		cin <= '1';

		wait for 3 ns;
		a <= '1';
		b <= '1';
		cin <= '0';

		wait for 3 ns;
		a <= '1';
		b <= '1';
		cin <= '1';
		end loop;
		wait;
	end process;
end;
