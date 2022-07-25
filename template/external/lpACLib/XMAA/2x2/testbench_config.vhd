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
-- Testbench for both accuracy configurable 2x2 multipliers
-- Author: Walaa El-Harouni  
--------------------------------------------
LIBRARY ieee;
USE ieee.std_logic_1164.ALL;
USE ieee.std_logic_unsigned.all;

library VITAL;
use VITAL.all;

ENTITY testbench IS
END testbench;

ARCHITECTURE behavior OF testbench IS
	
	COMPONENT Config2x2MultLit
	port (A   	 : in std_logic_vector(1 downto 0);
	      B 	 : in std_logic_vector(1 downto 0);
	      en	 : in std_logic;
	      Product    : out std_logic_vector(3 downto 0) );
	END COMPONENT;

	signal a: std_logic_vector (1 downto 0);
	signal b: std_logic_vector (1 downto 0);
	signal enable: std_logic;
	signal Product: std_logic_vector (3 downto 0);

BEGIN
	uut: Config2x2MultLit port map (A => a, B => b, en => enable, Product => Product);
		
	stim_proc : process
	begin
		for l in 0 to 50000 loop
		enable <= '0';

		wait for 2 ns;
		a <= "11";
		b <= "00";

		wait for 5 ns;
		a <= "11";
		b <= "01";

		wait for 5 ns;
		a <= "11";
		b <= "10";

		wait for 5 ns;
		a <= "11";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "00";
		b <= "00";

		wait for 5 ns;
		a <= "00";
		b <= "01";

		wait for 5 ns;
		a <= "00";
		b <= "10";

		wait for 5 ns;
		a <= "00";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "01";
		b <= "00";

		wait for 5 ns;
		a <= "01";
		b <= "01";

		wait for 5 ns;
		a <= "01";
		b <= "10";

		wait for 5 ns;
		a <= "01";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "10";
		b <= "00";

		wait for 5 ns;
		a <= "10";
		b <= "01";

		wait for 5 ns;
		a <= "10";
		b <= "10";

		wait for 5 ns;
		a <= "10";
		b <= "11";
		-- ========
		wait for 5 ns;

		enable <= '1';
		a <= "11";
		b <= "00";

		wait for 5 ns;
		a <= "11";
		b <= "01";

		wait for 5 ns;
		a <= "11";
		b <= "10";

		wait for 5 ns;
		a <= "11";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "00";
		b <= "00";

		wait for 5 ns;
		a <= "00";
		b <= "01";

		wait for 5 ns;
		a <= "00";
		b <= "10";

		wait for 5 ns;
		a <= "00";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "01";
		b <= "00";

		wait for 5 ns;
		a <= "01";
		b <= "01";

		wait for 5 ns;
		a <= "01";
		b <= "10";

		wait for 5 ns;
		a <= "01";
		b <= "11";
		-- ==========
		wait for 5 ns;
		a <= "10";
		b <= "00";

		wait for 5 ns;
		a <= "10";
		b <= "01";

		wait for 5 ns;
		a <= "10";
		b <= "10";

		wait for 5 ns;
		a <= "10";
		b <= "11";
		-- ========
		end loop;
		wait;
	end process;
end;
