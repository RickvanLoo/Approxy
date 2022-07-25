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
-- Testbench for all Multi-bit FA 
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
	
	COMPONENT AdderIMPACTFirstApproxMultiBit 
	--TODO: change to match settings in source file (vhdl), this is an example for 8-bit adder
	generic (bitWidth : integer := 8; approxBits : integer := 4); 
	port (A   	 : in std_logic_vector(bitWidth-1 downto 0);
		  B 	 : in std_logic_vector(bitWidth-1 downto 0);
		  Cin 	 : in std_logic;
		  Sum 	 : out std_logic_vector(bitWidth-1 downto 0);
		  Cout   : out std_logic );
	end component;

	signal a: std_logic_vector(7 downto 0);
	signal b: std_logic_vector(7 downto 0);
	signal cin: std_logic;
	signal sum: std_logic_vector(7 downto 0);
	signal cout: std_logic;

BEGIN
	uut : AdderIMPACTFirstApproxMultiBit port map (A => a, B => b, Cin => cin, Sum => sum, Cout => cout);
		
	— change according to number of bits
	stim_proc : process
	begin
		wait for 2 ns;
		for i in 0 to 6 loop
			for j in 0 to 255 loop
				a <= std_logic_vector(to_unsigned(j, 8));
				for k in 0 to 255 loop
					cin <= '0';
					b <= std_logic_vector(to_unsigned(k, 8));
					wait for 3 ns;
					cin <= '1';
					wait for 3 ns;
				end loop;
			end loop;
		end loop;
		wait;
	end process;
end;
