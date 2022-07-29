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
-- Config4x4MultV1First
-- Author: Walaa El-Harouni  
-- uses: Config2x2MultV1.vhd, Approx2x2MultV1.vhd, AdderIMPACTFirstApproxMultiBit.vhd, AdderAccurateOneBit.vhd and AdderIMPACTFirstApproxOneBit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity {{.EntityName}} is
generic (word_size: integer:={{.BitSize}}); 
Port ( 
A : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
B : in  STD_LOGIC_VECTOR (word_size-1 downto 0);
prod: out STD_LOGIC_VECTOR (word_size * 2 - 1 downto 0));
end {{.EntityName}};


architecture config4x4MultV1FirstArch of {{.EntityName}} is	
	component Config2x2MultV1 is
	port (A   	 : in std_logic_vector(1 downto 0);
		  B 	 : in std_logic_vector(1 downto 0);
		  en	 : in std_logic;
		  Product: out std_logic_vector(3 downto 0) );
	end component;
	
	component AdderIMPACTFirstApproxMultiBit is
	generic (bitWidth : integer := 8; approxBits : integer := 6);
	port (A   	 : in std_logic_vector(bitWidth-1 downto 0);
		  B 	 : in std_logic_vector(bitWidth-1 downto 0);
		  Cin 	 : in std_logic;
		  Sub 	 : in std_logic; -- '0' to add and '1' to subtract
		  Sum 	 : out std_logic_vector(bitWidth-1 downto 0);
		  Cout   : out std_logic );
	end component;
	
	signal LL, LH, HL, HH: std_logic_vector(3 downto 0);
	signal shifted_LL, shifted_LH, shifted_HL, shifted_HH: unsigned(7 downto 0);
	signal sum1, sum2, sum3 : std_logic_vector(7 downto 0);
	signal cout1, cout2, cout3: std_logic;
	
	constant TWO: natural := 2;
	constant FOUR: natural := 4;
	constant EIGHT: natural := 8;

	constant MAX1: natural := 3;
	constant MIN1: natural := 2;
	
	constant MAX2: natural := 1;
	constant MIN2: natural := 0;

	constant approxLSB: integer := 0;
	
	begin
	
	LxL: Config2x2MultV1 port map(A => A(MAX2 downto MIN2) , B => B(MAX2 downto MIN2) , en => '1', Product => LL);	
	HxL: Config2x2MultV1 port map(A => A(MAX1 downto MIN1) , B => B(MAX2 downto MIN2) , en => '1', Product => HL);
	LxH: Config2x2MultV1 port map(A => A(MAX2 downto MIN2) , B => B(MAX1 downto MIN1) , en => '1', Product => LH);
	HxH: Config2x2MultV1 port map(A => A(MAX1 downto MIN1) , B => B(MAX1 downto MIN1) , en => '1', Product => HH);
	
	-- shifting
	shifted_LL <= resize(unsigned(LL), EIGHT);
	shifted_LH <= shift_left( resize(unsigned(LH), EIGHT) , TWO);
	shifted_HL <= shift_left( resize(unsigned(HL), EIGHT) , TWO);
	shifted_HH <= shift_left( resize(unsigned(HH), EIGHT) , FOUR);
	
	-- addition 
	adder1: AdderIMPACTFirstApproxMultiBit generic map (approxBits => approxLSB) 
		port map (A => std_logic_vector(shifted_LL) , B =>  std_logic_vector(shifted_HL) , Cin => '0', Sub => '0', Sum => sum1, Cout => cout1);
			 
	adder2: AdderIMPACTFirstApproxMultiBit generic map (approxBits => approxLSB) 
		port map(A => sum1 , B =>  std_logic_vector(shifted_LH) , Cin => cout1, Sub => '0', Sum => sum2 , Cout => cout2);
			
	adder3: AdderIMPACTFirstApproxMultiBit generic map (approxBits => approxLSB) 
		port map(A => sum2 , B =>  std_logic_vector(shifted_HH) , Cin => cout2,  Sub => '0', Sum => sum3, Cout =>cout3);

	process(cout3, sum3)
	begin
		if(cout3 = '0') then
			prod <= std_logic_vector(resize(unsigned(sum3), EIGHT));
		else
			prod <= std_logic_vector(unsigned(sum3) + "10000000");  -- 128
		end if;
	end process;
	
	
end config4x4MultV1FirstArch;
