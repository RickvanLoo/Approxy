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
-- Approx2x2MultV4
-- Author: Walaa El-Harouni  
-- Implementation for approximate multiplier (Version 4) 
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;


entity Approx2x2MultV4 is
	port (A   	 : in std_logic_vector(1 downto 0);
		  B 	 : in std_logic_vector(1 downto 0);
		  Output    : out std_logic_vector(2 downto 0) );
end Approx2x2MultV4;

architecture approx2x2MultV4Arch of Approx2x2MultV4 is
	
begin
	process (A, B)
	begin 
		if A(1) = '0' and A(0) = '1' then
			Output(0) <= B(0);
			Output(1) <= B(1);
			Output(2) <= A(1);
		elsif A(1) = '1' and A(0) = '0' then
			Output(2) <= B(1);
			Output(1) <= B(0);
			Output(0) <= A(0);
		elsif A(1) = '1' and A(0) = '1' then
			if B(1) = '0' then
				Output(2) <= B(1);
				Output(1) <= B(0);
				Output(0) <= B(0);
			elsif B(1) = '1' then
				Output(2) <= not B(0);
				Output(1) <= not B(0);
				Output(0) <= B(0);
			else
				Output(2) <= '0';
				Output(1) <= '0';
				Output(0) <= '0';
			end if;
		else
			Output(2) <= '0';
			Output(1) <= '0';
			Output(0) <= '0';
		end if;
	end process;
	
end approx2x2MultV4Arch;
