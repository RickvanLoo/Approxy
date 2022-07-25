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
-- Approx4x4MultV1
-- Author: Walaa El-Harouni  
-- uses: Approx2x2MultV1.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;


entity Approx4x4MultV1 is
	port (A   	 : in std_logic_vector(3 downto 0);
		  B 	 : in std_logic_vector(3 downto 0);
		  Output    : out std_logic_vector(7 downto 0) );
end Approx4x4MultV1;

architecture approx4x4MultV1Arch of Approx4x4MultV1 is
	component Approx2x2MultV1 is
	port (A   	 : in std_logic_vector(1 downto 0);
	      B 	 : in std_logic_vector(1 downto 0);
	      Output     : out std_logic_vector(3 downto 0) );
	end component;
		
	signal LL, LH, HL, HH: std_logic_vector(3 downto 0);
	signal shifted_LL, shifted_LH, shifted_HL, shifted_HH: unsigned(7 downto 0);
	
	constant TWO: natural := 2;
	constant FOUR: natural := 4;
	constant EIGHT: natural := 8;

	constant MAX1: natural := 3;
	constant MIN1: natural := 2;
	
	constant MAX2: natural := 1;
	constant MIN2: natural := 0;
	
begin

	LxL: Approx2x2MultV1 port map(A => A(MAX2 downto MIN2) , B => B(MAX2 downto MIN2) , Output => LL);
	HxL: Approx2x2MultV1 port map(A => A(MAX1 downto MIN1) , B => B(MAX2 downto MIN2) , Output => HL);
	LxH: Approx2x2MultV1 port map(A => A(MAX2 downto MIN2) , B => B(MAX1 downto MIN1) , Output => LH);
	HxH: Approx2x2MultV1 port map(A => A(MAX1 downto MIN1) , B => B(MAX1 downto MIN1) , Output => HH);
	
	-- shifting
	shifted_LL <= resize(unsigned(LL), EIGHT);
	shifted_LH <= shift_left( resize(unsigned(LH), EIGHT) , TWO);
	shifted_HL <= shift_left( resize(unsigned(HL), EIGHT) , TWO);
	shifted_HH <= shift_left( resize(unsigned(HH), EIGHT) , FOUR);
		
	-- addition
	Output <= std_logic_vector(shifted_LL + shifted_LH + shifted_HL + shifted_HH);
		
end approx4x4MultV1Arch;




