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
-- Testbench for both accuracy configurable 4x4 multipliers
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
	
	COMPONENT Config4x4MultLit
		port (A  : in std_logic_vector(3 downto 0);
		      B  : in std_logic_vector(3 downto 0);
		     en	 : in std_logic;
		  Output : out std_logic_vector(7 downto 0) );
	END COMPONENT;

	signal a: std_logic_vector (3 downto 0);
	signal b: std_logic_vector (3 downto 0);
	signal enable: std_logic;
	signal Product: std_logic_vector (7 downto 0);

BEGIN
	uut: Config4x4MultLit port map (A => a, B => b, en => enable, Output => Product);
		
	stim_proc : process
	begin
		wait for 2 ns;
		for l in 0 to 500 loop
			for l in 0 to 15 loop
				enable <= '0';
				a <= std_logic_vector(to_unsigned(l, 4));

				b <= "0000"; wait for 5 ns;
				b <= "0001"; wait for 5 ns;
				b <= "0010"; wait for 5 ns;
				b <= "0011"; wait for 5 ns;
		
				b <= "0100"; wait for 5 ns;
				b <= "0101"; wait for 5 ns;
				b <= "0110"; wait for 5 ns;
				b <= "0111"; wait for 5 ns;

				b <= "1000"; wait for 5 ns;
				b <= "1001"; wait for 5 ns;
				b <= "1010"; wait for 5 ns;
				b <= "1011"; wait for 5 ns;

				b <= "1100"; wait for 5 ns;
				b <= "1101"; wait for 5 ns;
				b <= "1110"; wait for 5 ns;
				b <= "1111"; wait for 5 ns;

				enable <= '1';
				b <= "0000"; wait for 5 ns;
				b <= "0001"; wait for 5 ns;
				b <= "0010"; wait for 5 ns;
				b <= "0011"; wait for 5 ns;
		
				b <= "0100"; wait for 5 ns;
				b <= "0101"; wait for 5 ns;
				b <= "0110"; wait for 5 ns;
				b <= "0111"; wait for 5 ns;

				b <= "1000"; wait for 5 ns;
				b <= "1001"; wait for 5 ns;
				b <= "1010"; wait for 5 ns;
				b <= "1011"; wait for 5 ns;

				b <= "1100"; wait for 5 ns;
				b <= "1101"; wait for 5 ns;
				b <= "1110"; wait for 5 ns;
				b <= "1111"; wait for 5 ns;
			end loop;
		end loop;
		wait;
	end process;
end;
