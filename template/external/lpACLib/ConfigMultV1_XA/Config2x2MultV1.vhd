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
-- Approx2x2MultV4
-- Author: Walaa El-Harouni  
-- Implementation for an accuracy configurale approximate multiplier (Version 1)
-- uses: Approx2x2MultV1.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;


entity Config2x2MultV1 is
	port (A   	 : in std_logic_vector(1 downto 0);
		  B 	 : in std_logic_vector(1 downto 0);
		  en	 : in std_logic;
		  Product    : out std_logic_vector(3 downto 0) );
end Config2x2MultV1;

architecture config2x2MultV1Arch of Config2x2MultV1 is
	component Approx2x2MultV1 is
	port (A   	 : in std_logic_vector(1 downto 0);
		  B 	 : in std_logic_vector(1 downto 0);
		  Output : out std_logic_vector(3 downto 0) );
	end component;

	signal Output: std_logic_vector(3 downto 0);
	signal error: std_logic;        
begin
	inaccurate: Approx2x2MultV1 port map(A => A, B => B, Output => Output);
	
	error <= (Not Output(3)) and (A(0) and B(0));
	process (en, error, Output)
	begin
		if error='1' and en='1' then
			Product(0) <= not Output(0);
		else
			Product(0) <= Output(0);
		end if;	
			Product(3) <= Output(3);
			Product(2) <= Output(2);
			Product(1) <= Output(1);
	end process;
	
end config2x2MultV1Arch;
