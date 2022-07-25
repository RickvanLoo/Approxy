-- "lpACLib" is a library for Low-Power Approximate Computing Modules.
-- Copyright (C) 2016, Walaa El-Harouni, Muhammad Shafique, CES, KIT.
-- E-mail: walaa.elharouny@gmail.com, swahlah@yahoo.com

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
-- Testbench for accuracy configurable 4x4 multiplier from: 
-- "Trading Accuracy for Power with an Underdesigned Multiplier Architecture"
-- using AdderIMPACTThirdApproxMultiBit (IM_X5)
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
	
	COMPONENT ConfigThird4x4
	port (X   	       : in std_logic_vector(3 downto 0);
		  Y 	       : in std_logic_vector(3 downto 0);
		  accuracy     : in std_logic_vector(3 downto 0); -- size = #blocks  (4 for 4x4) [LxL, HxL, LxH, HxH]
		  Z            : out std_logic_vector(8 downto 0));
	END COMPONENT;
	
	
	signal a: std_logic_vector (3 downto 0);
	signal b: std_logic_vector (3 downto 0);
	signal enable: std_logic_vector (3 downto 0);
	signal Product: std_logic_vector (8 downto 0);

BEGIN
	uut: ConfigThird4x4 port map (X => a, Y => b, accuracy => enable, Z => Product);
		
	stim_proc : process
	begin
		wait for 2 ns;
		for l in 0 to 16 loop
			a <= std_logic_vector(to_unsigned(l, 4));

			for m in 0 to 16 loop
				enable <= std_logic_vector(to_unsigned(0, 4));
				b <= std_logic_vector(to_unsigned(m, 4));
				wait for 5 ns;

				enable <= std_logic_vector(to_unsigned(15, 4));
				wait for 5 ns;
	
			end loop;
		end loop;
		wait;
	end process;
end;
