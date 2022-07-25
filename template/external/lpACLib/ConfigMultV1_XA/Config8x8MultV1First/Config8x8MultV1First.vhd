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
-- Config8x8MultV1First
-- Author: Walaa El-Harouni  
-- uses: Config4x4MultV1First.vhd, Config2x2MultV1.vhd, Approx2x2MultV1.vhd,
--       AdderIMPACTFirstApproxMultiBit.vhd, AdderAccurateOneBit.vhd and AdderIMPACTFirstApproxOneBit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity Config8x8MultV1First is
	port (X   	       : in std_logic_vector(7 downto 0);
		  Y 	       : in std_logic_vector(7 downto 0);
		  accuracy     : in std_logic_vector(15 downto 0); -- size = #blocks  (16 for 4x4) [LxL, HxL, LxH, HxH]
		  Z            : out std_logic_vector(16 downto 0));
end Config8x8MultV1First;

architecture config8x8MultV1FirstArch of Config8x8MultV1First is	
	component Config4x4MultV1First is
	port (X   	       : in std_logic_vector(3 downto 0);
		  Y 	       : in std_logic_vector(3 downto 0);
		  accuracy     : in std_logic_vector(3 downto 0); -- size = #blocks  (4 for 4x4) [LxL, HxL, LxH, HxH]
		  Z            : out std_logic_vector(8 downto 0));
	end component;

	component AdderIMPACTFirstApproxMultiBit is
	generic (bitWidth : integer := 32; approxBits : integer := 6);
	port (A   	 : in std_logic_vector(bitWidth-1 downto 0);
		  B 	 : in std_logic_vector(bitWidth-1 downto 0);
		  Cin 	 : in std_logic;
		  Sub 	 : in std_logic; -- '0' to add and '1' to subtract
		  Sum 	 : out std_logic_vector(bitWidth-1 downto 0);
		  Cout   : out std_logic );
	end component;
	
	signal LL, LH, HL, HH: std_logic_vector(8 downto 0);
	signal shifted_LL, shifted_LH, shifted_HL, shifted_HH: unsigned(15 downto 0);
	signal sum1, sum2, sum3 : std_logic_vector(15 downto 0);
	signal cout1, cout2, cout3: std_logic;
	
	constant FOUR: natural := 4;
	constant EIGHT: natural := 8;
	constant SIXTEEN: natural := 16;

	constant MAX1: natural := 7;
	constant MIN1: natural := 4;
	
	constant MAX2: natural := 3;
	constant MIN2: natural := 0;
	
	constant approxLSB: integer := 8;
	
begin
	
	LxL: Config4x4MultV1First port map(X => X(MAX2 downto MIN2) , Y => Y(MAX2 downto MIN2) , accuracy => accuracy(3 downto 0), Z => LL);
	HxL: Config4x4MultV1First port map(X => X(MAX1 downto MIN1) , Y => Y(MAX2 downto MIN2) , accuracy => accuracy(7 downto 4), Z => HL);
	LxH: Config4x4MultV1First port map(X => X(MAX2 downto MIN2) , Y => Y(MAX1 downto MIN1) , accuracy => accuracy(11 downto 8), Z => LH);
	HxH: Config4x4MultV1First port map(X => X(MAX1 downto MIN1) , Y => Y(MAX1 downto MIN1) , accuracy => accuracy(15 downto 12), Z => HH);
	
	-- shifting
	shifted_LL <= resize(unsigned(LL), SIXTEEN);
	shifted_LH <= shift_left( resize(unsigned(LH), SIXTEEN) , FOUR);
	shifted_HL <= shift_left( resize(unsigned(HL), SIXTEEN) , FOUR);
	shifted_HH <= shift_left( resize(unsigned(HH), SIXTEEN) , EIGHT);
		
	-- addition
	adder1: AdderIMPACTFirstApproxMultiBit generic map (approxBits => approxLSB) 
		port map (A => std_logic_vector(shifted_LL) , B =>  std_logic_vector(shifted_HL) , Cin => '0', Sub => '0', Sum => sum1, Cout => cout1);
			 
	adder2: AdderIMPACTFirstApproxMultiBit generic map (approxBits => approxLSB) 
		port map(A => sum1 , B =>  std_logic_vector(shifted_LH) , Cin => cout1, Sub => '0', Sum => sum2 , Cout => cout2);
			
	adder3: AdderIMPACTFirstApproxMultiBit generic map (approxBits => approxLSB) 
		port map(A => sum2 , B =>  std_logic_vector(shifted_HH) , Cin => cout2, Sub => '0', Sum => sum3, Cout =>cout3);

	process(cout3, sum3)
	begin
		if(cout3 = '0') then
			Z <= std_logic_vector(resize(unsigned(sum3), SIXTEEN+1));
		elsif(cout3 = '1') then
			Z <= std_logic_vector(unsigned(sum3) + "10000000000000000");  -- 2^16
		end if;
	end process;
	
end config8x8MultV1FirstArch;
