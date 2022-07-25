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
-- Approx16x16MultLit
-- Author: Walaa El-Harouni  
-- uses: Approx8x8MultLit, Approx4x4MultLit.vhd and Approx2x2MultLit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;


entity Approx16x16MultLit is
	port (A   	 : in std_logic_vector(15 downto 0);
		  B 	 : in std_logic_vector(15 downto 0);
		  Output    : out std_logic_vector(31 downto 0) );
end Approx16x16MultLit;

architecture approx16x16MultLitArch of Approx16x16MultLit is
	component Approx8x8MultLit is
	port (A   	 : in std_logic_vector(7 downto 0);
		  B 	 : in std_logic_vector(7 downto 0);
		  Output    : out std_logic_vector(15 downto 0) );
	end component;
		
	signal LL, LH, HL, HH: std_logic_vector(15 downto 0);
	signal shifted_LL, shifted_LH, shifted_HL, shifted_HH: unsigned(31 downto 0);
	
	constant EIGHT: natural := 8;
	constant SIXTEEN: natural := 16;
	constant THIRTY_TWO: natural := 32;

	constant MAX1: natural := 15;
	constant MIN1: natural := 8;
	
	constant MAX2: natural := 7;
	constant MIN2: natural := 0;
	
begin

	LxL: Approx8x8MultLit port map(A => A(MAX2 downto MIN2) , B => B(MAX2 downto MIN2) , Output => LL);
	HxL: Approx8x8MultLit port map(A => A(MAX1 downto MIN1) , B => B(MAX2 downto MIN2) , Output => HL);
	LxH: Approx8x8MultLit port map(A => A(MAX2 downto MIN2) , B => B(MAX1 downto MIN1) , Output => LH);
	HxH: Approx8x8MultLit port map(A => A(MAX1 downto MIN1) , B => B(MAX1 downto MIN1) , Output => HH);
	
	-- shifting
	shifted_LL <= resize(unsigned(LL), THIRTY_TWO);
	shifted_LH <= shift_left( resize(unsigned(LH), THIRTY_TWO) , EIGHT);
	shifted_HL <= shift_left( resize(unsigned(HL), THIRTY_TWO) , EIGHT);
	shifted_HH <= shift_left( resize(unsigned(HH), THIRTY_TWO) , SIXTEEN);
		
	-- addition
	Output <= std_logic_vector(shifted_LL + shifted_LH + shifted_HL + shifted_HH);
		
end approx16x16MultLitArch;




