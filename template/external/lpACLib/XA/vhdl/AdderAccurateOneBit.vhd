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
-- AdderAccurateOneBit
-- Author: Walaa El-Harouni  
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity AdderAccurateOneBit is
	port (A   	 : in std_logic;
		  B 	 : in std_logic;
		  Cin 	 : in std_logic;
		  Sum 	 : out std_logic;
		  Cout   : out std_logic );
end AdderAccurateOneBit;

architecture AdderAccurateOneBitArch of AdderAccurateOneBit is
	signal axorb: std_logic;
begin
	axorb <= A xor B;
	
 	Sum <= axorb xor Cin;
	Cout <= (axorb and Cin) or (A and B);
	
end AdderAccurateOneBitArch;




