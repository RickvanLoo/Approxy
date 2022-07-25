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
-- Testbench for all 2x2 multipliers except the accuracy configurable ones
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
	
	-- change "Output" port to (2 downto 0) for Approx2x2MultLit, Approx2x2MultV2, Approx2x2MultV3, Approx2x2MultV4
	-- and "Output" signal also + initialise to "000" 
	COMPONENT Accurate2x2Mult 
	port (A: in std_logic_vector(1 downto 0); 
	      B: in std_logic_vector(1 downto 0);
 	      Output: out std_logic_vector(3 downto 0) ); 
	end component;

	signal a: std_logic_vector(1 downto 0);
	signal b: std_logic_vector(1 downto 0);
	signal Output: std_logic_vector(3 downto 0) := "0000";

BEGIN
	uut : Accurate2x2Mult port map (A => a, B => b, Output=>Output);
	
	stim_proc : process
	begin
		for l in 0 to 100000 loop
		wait for 2 ns;
		a <= "11";
		b <= "00";

		wait for 5 ns;
		a <= "11";
		b <= "01";

		wait for 5 ns;
		a <= "11";
		b <= "10";

		wait for 5 ns;
		a <= "11";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "00";
		b <= "00";

		wait for 5 ns;
		a <= "00";
		b <= "01";

		wait for 5 ns;
		a <= "00";
		b <= "10";

		wait for 5 ns;
		a <= "00";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "01";
		b <= "00";

		wait for 5 ns;
		a <= "01";
		b <= "01";

		wait for 5 ns;
		a <= "01";
		b <= "10";

		wait for 5 ns;
		a <= "01";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "10";
		b <= "00";

		wait for 5 ns;
		a <= "10";
		b <= "01";

		wait for 5 ns;
		a <= "10";
		b <= "10";

		wait for 5 ns;
		a <= "10";
		b <= "11";
		end loop;
		wait;
	end process;
end;
