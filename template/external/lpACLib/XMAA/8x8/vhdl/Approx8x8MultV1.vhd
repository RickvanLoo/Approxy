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
-- Approx8x8MultV1
-- Author: Walaa El-Harouni  
-- uses: Approx4x4MultV1.vhd and Approx2x2MultV1.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity Approx8x8MultV1 is
	port (A   	 : in std_logic_vector(7 downto 0);
		  B 	 : in std_logic_vector(7 downto 0);
		  Output    : out std_logic_vector(15 downto 0) );
end Approx8x8MultV1;

architecture approx8x8MultV1Arch of Approx8x8MultV1 is
	component Approx4x4MultV1 is
	port (A   	 : in std_logic_vector(3 downto 0);
	      B 	 : in std_logic_vector(3 downto 0);
	      Output     : out std_logic_vector(7 downto 0) );
	end component;
		
	signal LL, LH, HL, HH: std_logic_vector(7 downto 0);
	signal shifted_LL, shifted_LH, shifted_HL, shifted_HH: unsigned(15 downto 0);
	
	constant FOUR: natural := 4;
	constant EIGHT: natural := 8;
	constant SIXTEEN: natural := 16;

	constant MAX1: natural := 7;
	constant MIN1: natural := 4;
	
	constant MAX2: natural := 3;
	constant MIN2: natural := 0;
	
begin

	LxL: Approx4x4MultV1 port map(A => A(MAX2 downto MIN2) , B => B(MAX2 downto MIN2) , Output => LL);
	HxL: Approx4x4MultV1 port map(A => A(MAX1 downto MIN1) , B => B(MAX2 downto MIN2) , Output => HL);
	LxH: Approx4x4MultV1 port map(A => A(MAX2 downto MIN2) , B => B(MAX1 downto MIN1) , Output => LH);
	HxH: Approx4x4MultV1 port map(A => A(MAX1 downto MIN1) , B => B(MAX1 downto MIN1) , Output => HH);
	
	-- shifting
	shifted_LL <= resize(unsigned(LL), SIXTEEN);
	shifted_LH <= shift_left( resize(unsigned(LH), SIXTEEN) , FOUR);
	shifted_HL <= shift_left( resize(unsigned(HL), SIXTEEN) , FOUR);
	shifted_HH <= shift_left( resize(unsigned(HH), SIXTEEN) , EIGHT);
		
	-- addition
	Output <= std_logic_vector(shifted_LL + shifted_LH + shifted_HL + shifted_HH);
		
end approx8x8MultV1Arch;




