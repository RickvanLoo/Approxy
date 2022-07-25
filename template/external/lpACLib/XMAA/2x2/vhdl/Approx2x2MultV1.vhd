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
-- Approx2x2MultV1
-- Author: Walaa El-Harouni  
-- Implementation for approximate multiplier (Version 1) 
--------------------------------------------


library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity Approx2x2MultV1 is
	port (A   	 : in std_logic_vector(1 downto 0);
		  B 	 : in std_logic_vector(1 downto 0);
		  Output    : out std_logic_vector(3 downto 0) );
end Approx2x2MultV1;

architecture approx2x2MultV1Arch of Approx2x2MultV1 is

	signal a1b1, a1b0, a0b1, a0a1b0b1: std_logic;
    
begin
	a1b1 <= A(1) and B(1);
	a1b0 <= A(1) and B(0);
	a0b1 <= A(0) and B(1);
	
	a0a1b0b1 <= a0b1 and a1b0;
	
	Output(3) <= a0a1b0b1;
	Output(2) <= a1b1 xor a0a1b0b1;
	Output(1) <= a0b1 xor a1b0;
	Output(0) <= a0a1b0b1;
	
end approx2x2MultV1Arch;