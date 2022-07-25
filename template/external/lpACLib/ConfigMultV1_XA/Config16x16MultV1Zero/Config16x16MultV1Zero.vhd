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
-- Config16x16MultV1Zero
-- Author: Walaa El-Harouni  
-- uses: Config8x8MultV1Zero.vhd, Config4x4MultV1Zero.vhd, Config2x2MultV1.vhd, Approx2x2MultV1.vhd,
--       AdderIMPACTZeroApproxMultiBit.vhd, AdderAccurateOneBit.vhd and AdderIMPACTZeroApproxOneBit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity Config16x16MultV1Zero is
	port (X   	       : in std_logic_vector(15 downto 0);
		  Y 	       : in std_logic_vector(15 downto 0);
		  accuracy     : in std_logic_vector(63 downto 0); -- size = #blocks  (64 for 16x16) [LxL, HxL, LxH, HxH]
		  Z            : out std_logic_vector(31 downto 0));
end Config16x16MultV1Zero;

architecture config16x16MultV1ZeroArch of Config16x16MultV1Zero is	
	component Config8x8MultV1Zero is
	port (X   	       : in std_logic_vector(7 downto 0);
		  Y 	       : in std_logic_vector(7 downto 0);
		  accuracy     : in std_logic_vector(15 downto 0); -- size = #blocks  (16 for 8x8) [LxL, HxL, LxH, HxH]
		  Z            : out std_logic_vector(15 downto 0));
	end component;

	component AdderIMPACTZeroApproxMultiBit is
	generic (bitWidth : integer := 16; approxBits : integer := 0);
	port (A   	 : in std_logic_vector(bitWidth-1 downto 0);
		  B 	 : in std_logic_vector(bitWidth-1 downto 0);
		  Cin 	 : in std_logic;
		  Sum 	 : out std_logic_vector(bitWidth-1 downto 0);
		  Cout   : out std_logic );
	end component;
	

	signal sum1, sum2, sum3 : std_logic_vector(31 downto 0);
	signal cout1, cout2, cout3: std_logic;	
	signal LL, LH, HL, HH: std_logic_vector(15 downto 0);
	signal shifted_LL, shifted_LH, shifted_HL, shifted_HH: unsigned(31 downto 0);
	
	constant EIGHT: natural := 8;
	constant SIXTEEN: natural := 16;
	constant THIRTY_TWO: natural := 32;

	constant MAX1: natural := 15;
	constant MIN1: natural := 8;
	
	constant MAX2: natural := 7;
	constant MIN2: natural := 0;
	
	constant approxLSB: integer := 0;
	
begin
	
	LxL: Config8x8MultV1Zero port map(X => X(MAX2 downto MIN2) , Y => Y(MAX2 downto MIN2) , accuracy => accuracy(15 downto 0), Z => LL);
	HxL: Config8x8MultV1Zero port map(X => X(MAX1 downto MIN1) , Y => Y(MAX2 downto MIN2) , accuracy => accuracy(31 downto 16), Z => HL);
	LxH: Config8x8MultV1Zero port map(X => X(MAX2 downto MIN2) , Y => Y(MAX1 downto MIN1) , accuracy => accuracy(47 downto 32), Z => LH);
	HxH: Config8x8MultV1Zero port map(X => X(MAX1 downto MIN1) , Y => Y(MAX1 downto MIN1) , accuracy => accuracy(63 downto 48), Z => HH);
	
	-- shifting
	shifted_LL <= resize(unsigned(LL), THIRTY_TWO);
	shifted_LH <= shift_left( resize(unsigned(LH), THIRTY_TWO) , EIGHT);
	shifted_HL <= shift_left( resize(unsigned(HL), THIRTY_TWO) , EIGHT);
	shifted_HH <= shift_left( resize(unsigned(HH), THIRTY_TWO) , SIXTEEN);
		
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
			Z <= sum3;
		else
			Z <= std_logic_vector(unsigned(sum3) + "10000000000000000000000000000000");  -- (2^31)
		end if;
	end process;
	
end config16x16MultV1ZeroArch;
