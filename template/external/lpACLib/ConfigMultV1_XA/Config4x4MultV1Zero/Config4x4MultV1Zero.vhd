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
-- Config4x4MultV1Zero
-- Author: Walaa El-Harouni  
-- uses: Config2x2MultV1.vhd, Approx2x2MultV1.vhd, AdderIMPACTZeroApproxMultiBit.vhd, AdderAccurateOneBit.vhd and AdderIMPACTZeroApproxOneBit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity Config4x4MultV1Zero is
	port (X   	       : in std_logic_vector(3 downto 0);
		  Y 	       : in std_logic_vector(3 downto 0);
		  accuracy     : in std_logic_vector(3 downto 0); -- size = #blocks  (4 for 4x4) [LxL, HxL, LxH, HxH]
		  Z            : out std_logic_vector(8 downto 0));
end Config4x4MultV1Zero;

architecture config4x4MultV1ZeroArch of Config4x4MultV1Zero is	
	component Config2x2MultV1 is
	port (A   	 : in std_logic_vector(1 downto 0);
		  B 	 : in std_logic_vector(1 downto 0);
		  en	 : in std_logic;
		  Product: out std_logic_vector(3 downto 0) );
	end component;
	
	component AdderIMPACTZeroApproxMultiBit is
	generic (bitWidth : integer := 32; approxBits : integer := 6);
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

	constant approxLSB: integer := 8;
	
	begin
	
	LxL: Config2x2MultV1 port map(A => X(MAX2 downto MIN2) , B => Y(MAX2 downto MIN2) , en => accuracy(0), Product => LL);	
	HxL: Config2x2MultV1 port map(A => X(MAX1 downto MIN1) , B => Y(MAX2 downto MIN2) , en => accuracy(1), Product => HL);
	LxH: Config2x2MultV1 port map(A => X(MAX2 downto MIN2) , B => Y(MAX1 downto MIN1) , en => accuracy(2), Product => LH);
	HxH: Config2x2MultV1 port map(A => X(MAX1 downto MIN1) , B => Y(MAX1 downto MIN1) , en => accuracy(3), Product => HH);
	
	-- shifting
	shifted_LL <= resize(unsigned(LL), EIGHT);
	shifted_LH <= shift_left( resize(unsigned(LH), EIGHT) , TWO);
	shifted_HL <= shift_left( resize(unsigned(HL), EIGHT) , TWO);
	shifted_HH <= shift_left( resize(unsigned(HH), EIGHT) , FOUR);
	
	-- addition 
	adder1: AdderIMPACTZeroApproxMultiBit generic map (approxBits => approxLSB) 
		port map (A => std_logic_vector(shifted_LL) , B =>  std_logic_vector(shifted_HL) , Cin => '0', Sub => '0', Sum => sum1, Cout => cout1);
			 
	adder2: AdderIMPACTZeroApproxMultiBit generic map (approxBits => approxLSB) 
		port map(A => sum1 , B =>  std_logic_vector(shifted_LH) , Cin => cout1, Sub => '0', Sum => sum2 , Cout => cout2);
			
	adder3: AdderIMPACTZeroApproxMultiBit generic map (approxBits => approxLSB) 
		port map(A => sum2 , B =>  std_logic_vector(shifted_HH) , Cin => cout2, Sub => '0', Sum => sum3, Cout =>cout3);

	process(cout3, sum3)
	begin
		if(cout3 = '0') then
			Z <= std_logic_vector(resize(unsigned(sum3), EIGHT+1));
		elsif(cout3 = '1') then
			Z <= std_logic_vector(unsigned(sum3) + "100000000");  -- 2^8
		end if;
	end process;
	
	
end config4x4MultV1ZeroArch;
