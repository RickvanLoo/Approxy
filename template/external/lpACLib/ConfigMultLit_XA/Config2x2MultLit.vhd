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
-- Config2x2MultLit
-- Author: Walaa El-Harouni  
-- Implementation for an accuracy configurale approximate multiplier from:
-- "Trading Accuracy for Power with an Underdesigned Multiplier Architecture" 
-- uses: Approx2x2MultLit.vhd
--------------------------------------------

library ieee;
use ieee.std_logic_1164.ALL;
use ieee.numeric_std.ALL;

entity Config2x2MultLit is
	port (A   	 : in std_logic_vector(1 downto 0);
		  B 	 : in std_logic_vector(1 downto 0);
		  en	 : in std_logic;
		  Product    : out std_logic_vector(3 downto 0) );
end Config2x2MultLit;

architecture config2x2MultLitArch of Config2x2MultLit is
	component Approx2x2MultLit is
	port (A   	 : in std_logic_vector(1 downto 0);
	      B 	 : in std_logic_vector(1 downto 0);
	      Output     : out std_logic_vector(2 downto 0) );
	end component;
	    
    signal approxOut: std_logic_vector(2 downto 0);
    
begin
	inaccurate: Approx2x2MultLit port map(A => A, B => B, Output => approxOut);

	process (en, approxOut)

    	variable inaccurateOut: integer;
    	variable result: integer := 0;

	begin
		if en = '1' and approxOut = "111" then
			inaccurateOut := to_integer(unsigned(approxOut));
			result := 2 + inaccurateOut; 
			Product <= std_logic_vector(to_unsigned(result, 4));
		else
			Product <= std_logic_vector(resize(unsigned(approxOut), 4));
		end if;	
		
	end process;
	
end config2x2MultLitArch;
