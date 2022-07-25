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
-- Testbench for Version 1 accuracy configurable 4x4 multiplier using AdderIMPACTFirstApproxMultiBit 
-- Author: Walaa El-Harouni  
--------------------------------------------

LIBRARY ieee;
USE ieee.std_logic_1164.ALL;
USE ieee.std_logic_unsigned.all;
use ieee.numeric_std.ALL;

library VITAL;
use VITAL.all;

ENTITY testbench IS
END testbench;

ARCHITECTURE behavior OF testbench IS
	
	COMPONENT Config4x4MultV1First
		port (X   	   : in std_logic_vector(3 downto 0);
		      Y 	   : in std_logic_vector(3 downto 0);
		      accuracy     : in std_logic_vector(3 downto 0); -- size = #blocks  (4 for 4x4) [LxL, HxL, LxH, HxH]
		      Z            : out std_logic_vector(7 downto 0));
	END COMPONENT;

	signal a: std_logic_vector (3 downto 0);
	signal b: std_logic_vector (3 downto 0);
	signal enable: std_logic_vector(3 downto 0);
	signal Product: std_logic_vector (7 downto 0);

BEGIN
	uut: Config4x4MultV1First port map (X => a, Y => b, accuracy => enable, Z => Product);
		
	stim_proc : process
	begin
		enable <= "1111";
		wait for 2 ns;

		for l in 0 to 1000 loop
			for l in 0 to 15 loop
				a <= std_logic_vector(to_unsigned(l, 4));
				for k in 0 to 15 loop
					b <= std_logic_vector(to_unsigned(k, 4));
				        wait for 5 ns;
				end loop;
			end loop;
		end loop;
		wait;
	end process;
end;
